from sqlalchemy.ext.asyncio import create_async_engine, async_sessionmaker, AsyncSession
from sqlalchemy.orm import DeclarativeBase
from core.config import settings

# CHUẨN: Gọi từ instance (s viết thường)
engine = create_async_engine(settings.DATA_TABLE_URL, echo=True, pool_pre_ping=True)

async_session = async_sessionmaker(bind=engine, class_=AsyncSession, expire_on_commit=False)

class Base(DeclarativeBase):
    pass

async def get_db_session():
    """
    Generation này
    :return:
    """
    async with async_session() as session:
        try:
            yield session
        finally:
            await session.close()