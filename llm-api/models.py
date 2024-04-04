from pydantic import BaseModel, field_validator


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


class TextChunk(BaseModel):
    id: str
    chunk: str
