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


def init_mock_response() -> bool:
    return os.environ.get('MOCK_RESPONSE', 'false') == 'true'


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
API_KEY = os.environ.get('API_KEY', '')
client = AsyncClient(
    timeout=60,
)
PROVIDER: providers.Provider = (
    providers.Ollama(client, BASE_URL, MODEL)
    if get_provider() == providers.Ollama.NAME
    else providers.OpenAI(client, BASE_URL, API_KEY, MODEL)
    if get_provider() == providers.OpenAI.NAME
    else providers.Google(client, API_KEY, MODEL)
)

MOCK_RESPONSE = init_mock_response()


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
    if MOCK_RESPONSE:
        return MulipleChoice(
            question='What is the capital of France?',
            options=['Paris', 'London', 'Berlin', 'Madrid'],
            correct_options=['A'],
        )
    lang = detect_lang.detect_language(context.prompt)
    messages = llmio.format_question(
        context.prompt, models.MULTIPLE_CHOICE, lang
    )
    response = await PROVIDER.get_model_response(messages)
    question = llmio.try_parse_multiple_choice(response)
    if question is None:
        raise ValueError('Failed to generate single choice question')
    return question


@app.post('/single-choice/create')
async def single_choice_create(context: Prompt) -> SingleChoice:
    if MOCK_RESPONSE:
        return SingleChoice(
            question='What is the capital of France?',
            options=['Paris', 'London', 'Berlin', 'Madrid'],
            correct_option='A',
        )
    lang = detect_lang.detect_language(context.prompt)
    messages = llmio.format_question(
        context.prompt, models.SINGLE_CHOICE, lang
    )
    response = await PROVIDER.get_model_response(messages)
    question = llmio.try_parse_single_choice(response)
    if question is None:
        raise ValueError('Failed to generate single choice question')
    return question


@app.post('/true-or-false/create')
async def true_or_false_create(context: Prompt) -> TrueOrFalse:
    if MOCK_RESPONSE:
        return TrueOrFalse(
            question='Budapest is the capital of Hungary.', correct_option=True
        )
    lang = detect_lang.detect_language(context.prompt)
    messages = llmio.format_question(
        context.prompt, models.TRUE_OR_FALSE, lang
    )
    response = await PROVIDER.get_model_response(messages)
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
