from pydantic import BaseModel, field_validator


class MulipleChoice(BaseModel):
    question: str
    options: list[str]
    correct_options: list[str]

    @field_validator('options', mode='before')
    @classmethod
    def valid_options(cls, v):
        if len(v) != 4:
            raise ValueError('options must be a list of 4 strings')
        return v

    @field_validator('correct_options', mode='before')
    @classmethod
    def valid_correct_option(cls, v):
        if len(v) > 4 or len(v) < 1:
            raise ValueError('correct_option must be 1 or 4 characters')
        for c in v:
            if c not in 'ABCD':
                raise ValueError('correct_options must be [A, B, C, D]')
        return v


class SingleChoice(BaseModel):
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
        if v not in 'ABCD':
            raise ValueError('correct_option must be [A, B, C, D]')
        if len(v) > 1 or len(v) == 0:
            raise ValueError('correct_option must be 1 character')
        return v


class TrueOrFalse(BaseModel):
    question: str
    correct_option: bool


class Prompt(BaseModel):
    prompt: str


class TextChunk(BaseModel):
    id: str
    chunk: str
