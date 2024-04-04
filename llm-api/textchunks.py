import asyncpg
import uuid
from models import TextChunk
from returns.result import Result, Success, Failure
from langchain_text_splitters import CharacterTextSplitter

text_splitter = CharacterTextSplitter(
    separator='\n\n',
    chunk_size=1000,
    chunk_overlap=200,
    length_function=len,
    is_separator_regex=False,
)


def get_chunks(text: str) -> list[TextChunk]:
    return [
        TextChunk(id=str(uuid.uuid4()), chunk=chunk)
        for chunk in text_splitter.split_text(text)
    ]


async def init_chunk_storage(conn: asyncpg.Connection) -> Result[None, str]:
    schema = """
    CREATE TABLE IF NOT EXISTS textchunks (
        id UUID PRIMARY KEY,
        chunk TEXT NOT NULL
    );
    """
    try:
        await conn.execute(schema)
        return Success(None)
    except Exception as e:
        return Failure(f'Failed to create textchunks table: {e}')


async def insert_chunks(
    conn: asyncpg.Connection, chunks: list[TextChunk]
) -> Result[None, str]:
    query = """
    INSERT INTO textchunks (id, chunk)
    VALUES ($1, $2);
    """
    try:
        await conn.executemany(
            query, [(chunk.id, chunk.chunk) for chunk in chunks]
        )
        return Success(None)
    except Exception as e:
        print(e)
        return Failure(f'Failed to insert chunks: {e}')


async def get_chunk(
    conn: asyncpg.Connection, chunk_id: int
) -> Result[str, str]:
    query = """
    SELECT chunk FROM textchunks WHERE id = $1;
    """
    try:
        row = await conn.fetchrow(query, str(chunk_id))
        if row is None:
            return Failure(f'Chunk {chunk_id} not found')
        return Success(row['chunk'])
    except Exception as e:
        print(e)
        return Failure(f'Failed to get chunk: {e}')


async def delete_chunk(
    conn: asyncpg.Connection, chunk_id: int
) -> Result[None, str]:
    query = """
    DELETE FROM textchunks WHERE id = $1;
    """
    try:
        await conn.execute(query, str(chunk_id))
    except Exception as e:
        print(e)
        return Failure(f'Failed to delete chunk: {e}')
    return Success(None)
