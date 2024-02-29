import os
import json
from pydantic import BaseModel, ValidationError, field_validator
from fastapi import FastAPI
from httpx import AsyncClient

app = FastAPI()
client = AsyncClient()


def init_model() -> str:
    return os.environ.get('MODEL', 'mistral-quizgen')


def init_base_url() -> str:
    return os.environ.get('OLLAMA_URL', 'http://172.19.0.2:11434')


MODEL = init_model()
BASE_URL = init_base_url()


def ollama_request_data(prompt: str):
    return {
        'model': MODEL,
        'prompt': prompt,
        'stream': False,
        'format': 'json',
    }


class MulipleChoice(BaseModel):
    question: str
    options: list[str]
    correct_option: str

    @field_validator('options', mode='before')
    @classmethod
    def valid_options(cls, v):
        if len(v) != 4:
            raise ValueError('options must be a list of 4 strings')
        return v

    @field_validator('correct_option', mode='before')
    @classmethod
    def valid_correct_option(cls, v):
        if len(v) != 1:
            raise ValueError('correct_option must be a single character')
        if v not in 'ABCD':
            raise ValueError('correct_option must be one of A, B, C, D')
        return v


class Prompt(BaseModel):
    prompt: str


def try_parse_multiple_choice(data: str) -> MulipleChoice | None:
    try:
        response = json.loads(data)
        return MulipleChoice(
            question=response['Question'],
            options=response['Answers'],
            correct_option=response['Solution'],
        )
    except json.JSONDecodeError or ValidationError or KeyError:
        return None


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
