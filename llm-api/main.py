import os

import llmio
import database
import models
import providers
import textchunks
import detect_lang
from fastapi import FastAPI, HTTPException
from httpx import AsyncClient
from models import Prompt, MulipleChoice, SingleChoice, TextChunk, TrueOrFalse
from returns.pipeline import is_successful

app = FastAPI()


def get_provider() -> str:
    return os.environ.get('PROVIDER', providers.Ollama.NAME)


def init_model() -> str:
    provider = get_provider()
    if provider == providers.Ollama.NAME:
        return os.environ.get('MODEL', 'llama3.1:8b')
    else:
        return os.environ.get(
            'MODEL', 'jazzysnake01/llama-3-8b-quizgen-instruct'
        )


def init_base_url() -> str:
    return os.environ.get('BASE_URL', 'http://ollama:11434')


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
API_KEY = os.environ.get('API_KEY')
PROVIDER: providers.Provider = (
    providers.Ollama()
    if get_provider() == providers.Ollama.NAME
    else providers.OpenAI()
)

client = AsyncClient(
    base_url=BASE_URL,
    timeout=60 if get_provider() == providers.Ollama.NAME else 30,
    headers={'Authorization': f'Bearer {API_KEY}'}
    if API_KEY is not None
    else None,
)


def request_data(messages) -> dict:
    return {
        'temperature': 0.4,
        'model': MODEL,
        'stream': False,
        'messages': messages,
        'max_new_tokens': 500,
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
    lang = detect_lang.detect_language(context.prompt)
    messages = llmio.format_question(
        context.prompt, models.MULTIPLE_CHOICE, lang
    )
    res = await client.post(
        PROVIDER.CHAT_COMPLETION_ENDPOINT,
        json=request_data(messages),
    )
    res.raise_for_status()
    response = PROVIDER.parse_chat_completion_response(res.json())
    question = llmio.try_parse_multiple_choice(response)
    if question is None:
        raise ValueError('Failed to generate single choice question')
    return question


@app.post('/single-choice/create')
async def single_choice_create(context: Prompt) -> SingleChoice:
    lang = detect_lang.detect_language(context.prompt)
    messages = llmio.format_question(
        context.prompt, models.SINGLE_CHOICE, lang
    )
    res = await client.post(
        PROVIDER.CHAT_COMPLETION_ENDPOINT,
        json=request_data(messages),
    )
    res.raise_for_status()
    response = PROVIDER.parse_chat_completion_response(res.json())
    question = llmio.try_parse_single_choice(response)
    if question is None:
        raise ValueError('Failed to generate single choice question')
    return question


@app.post('/true-or-false/create')
async def true_or_false_create(context: Prompt) -> TrueOrFalse:
    lang = detect_lang.detect_language(context.prompt)
    messages = llmio.format_question(
        context.prompt, models.TRUE_OR_FALSE, lang
    )
    res = await client.post(
        PROVIDER.CHAT_COMPLETION_ENDPOINT,
        json=request_data(messages),
    )
    res.raise_for_status()
    response = PROVIDER.parse_chat_completion_response(res.json())
    question = llmio.try_parse_true_or_false(response)
    if question is None:
        raise ValueError('Failed to generate single choice question')
    return question


@app.post('/chunk')
async def chunk_create(context: Prompt) -> list[TextChunk]:
    chunks = textchunks.get_chunks(context.prompt)
    async with database.get_connection() as conn:
        res = await textchunks.insert_chunks(conn, chunks)
        if not is_successful(res):
            raise HTTPException(status_code=500, detail=res.failure())
    return chunks
