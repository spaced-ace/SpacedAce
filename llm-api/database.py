import asyncpg

from pgvector.asyncpg import register_vector
from contextlib import asynccontextmanager
from typing import Awaitable, Callable, List
from returns.pipeline import is_successful

from returns.result import Result


_DB_POOL = None


async def init(
    config: dict,
    repo_initializers: List[
        Callable[[asyncpg.Connection], Awaitable[Result[None, str]]]
    ],
):
    global _DB_POOL
    if _DB_POOL is None:
        _DB_POOL = await asyncpg.create_pool(**config)

    async with _DB_POOL.acquire() as conn:
        await conn.execute('CREATE EXTENSION IF NOT EXISTS vector')
        for initializer in repo_initializers:
            res = await initializer(conn)
            if not is_successful(res):
                print(res.failure())
                exit(1)


@asynccontextmanager
async def get_connection():
    global _DB_POOL
    if _DB_POOL is None:
        raise NameError('The module was not initialized')
    conn = await _DB_POOL.acquire()
    await register_vector(conn)
    try:
        yield conn
    finally:
        await _DB_POOL.release(conn)


async def close_pool():
    global _DB_POOL
    if _DB_POOL is not None:
        await _DB_POOL.close()
