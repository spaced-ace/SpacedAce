import os
import database
import textchunks
from fastapi import FastAPI, HTTPException
from httpx import AsyncClient
from models import Prompt, MulipleChoice, TextChunk
from multiple_choice import try_parse_multiple_choice
from returns.pipeline import is_successful

app = FastAPI()
client = AsyncClient()


def init_model() -> str:
    return os.environ.get('MODEL', 'mistral-quizgen')


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


def ollama_request_data(prompt: str):
    return {
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
        json=ollama_request_data(context.prompt),
        timeout=60,
    )
    response = res.json()['response']
    question = try_parse_multiple_choice(response)
    if question is None:
        raise ValueError('Failed to generate multiple choice question')
    return question


@app.post('/chunk')
async def chunk_create(context: Prompt) -> list[TextChunk]:
    chunks = textchunks.get_chunks(context.prompt)
    async with database.get_connection() as conn:
        res = await textchunks.insert_chunks(conn, chunks)
        if not is_successful(res):
            raise HTTPException(status_code=500, detail=res.failure())
    return chunks
