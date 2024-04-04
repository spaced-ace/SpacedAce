import os
import database
import textchunks
import multiple_choice
import single_choice
import true_or_false
from fastapi import FastAPI, HTTPException
from httpx import AsyncClient
from models import Prompt, MulipleChoice, SingleChoice, TextChunk, TrueOrFalse
from returns.pipeline import is_successful

app = FastAPI()
client = AsyncClient()


def init_model() -> str:
    return os.environ.get('MODEL', 'mistral')


def init_base_url() -> str:
    return os.environ.get('OLLAMA_URL', 'http://ollama:11434')


def get_db_config() -> dict:
    return {
        'host': os.environ.get('DB_HOST', 'database'),
        'port': os.environ.get('DB_PORT', '5432'),
        'user': os.environ.get('DB_USER', 'postgres'),
        'password': os.environ.get('DB_PASS', 'pw'),
        'database': os.environ.get('DB_NAME', 'postgres'),
    }


MODEL = init_model()
BASE_URL = init_base_url()


def ollama_request_data(prompt: str, template: str, system: str) -> dict:
    return {
        'template': template,
        'system': system,
        'temperature': 0.1,
        'model': MODEL,
        'prompt': prompt,
        'stream': False,
        'format': 'json',
    }


@app.on_event('startup')
async def startup_event():
    db_config = get_db_config()
    await database.init(
        db_config,
        [textchunks.init_chunk_storage],
    )


@app.on_event('shutdown')
async def close_resources():
    await database.close_pool()


@app.post('/multiple-choice/create')
async def multiple_choice_create(context: Prompt) -> MulipleChoice:
    res = await client.post(
        f'{BASE_URL}/api/generate',
        json=ollama_request_data(
            context.prompt, multiple_choice.TEMPLATE, multiple_choice.SYSTEM
        ),
        timeout=60,
    )
    response = res.json()['response']
    question = multiple_choice.try_parse_multiple_choice(response)
    if question is None:
        raise ValueError('Failed to generate multiple choice question')
    return question


@app.post('/single-choice/create')
async def single_choice_create(context: Prompt) -> SingleChoice:
    res = await client.post(
        f'{BASE_URL}/api/generate',
        json=ollama_request_data(
            context.prompt, single_choice.TEMPLATE, single_choice.SYSTEM
        ),
        timeout=60,
    )
    response = res.json()['response']
    question = single_choice.try_parse_single_choice(response)
    if question is None:
        raise ValueError('Failed to generate single choice question')
    return question


@app.post('/true-or-false/create')
async def true_or_false_create(context: Prompt) -> TrueOrFalse:
    res = await client.post(
        f'{BASE_URL}/api/generate',
        json=ollama_request_data(
            context.prompt, true_or_false.TEMPLATE, true_or_false.SYSTEM
        ),
        timeout=60,
    )
    response = res.json()['response']
    question = true_or_false.try_parse_true_or_false(response)
    if question is None:
        raise ValueError('Failed to generate true-or-false question')
    return question


@app.post('/chunk')
async def chunk_create(context: Prompt) -> list[TextChunk]:
    chunks = textchunks.get_chunks(context.prompt)
    async with database.get_connection() as conn:
        res = await textchunks.insert_chunks(conn, chunks)
        if not is_successful(res):
            raise HTTPException(status_code=500, detail=res.failure())
    return chunks
