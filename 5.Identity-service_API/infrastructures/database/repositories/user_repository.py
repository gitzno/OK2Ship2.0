from typing import Optional
from uuid import UUID

from fastapi.params import Security
from sqlalchemy import delete, insert, update
from sqlalchemy.future import select
from sqlalchemy.ext.asyncio import AsyncSession, result
from sqlalchemy.orm import class_mapper

from domain.interfaces.repositories.i_user_repository import IUserRepository
from domain.models import generated_models
from domain.models.generated_models import Users
from sqlalchemy.exc import IntegrityError

from domain.schemas.exceptions import DuplicateAccountError


class UserRepository(IUserRepository):

    def __init__(self, session: AsyncSession):
        self.session = session

    async def update_stamp(self, user_id: UUID, stamp: UUID) -> bool:
        query = (
            update(Users)
            .where(Users.UserID == user_id)
            .values(SecurityStamp=stamp)
        )

        result = await self.session.execute(query)

        return result.rowcount == 1

    async def get_by_id(self, user_id: UUID) -> Optional[Users]:
        stmt = select(Users).where(Users.UserID == user_id)
        result = await self.session.execute(stmt)
        return result.scalars().first()

    async def get_user_by_account(self, account: str) -> Optional[Users]:
        stmt = select(Users).filter(Users.Username == account)

        result = await self.session.execute(stmt)

        return result.scalars().first()

    async def create_user(self, user: Users) -> Users:
        try:
            user_dict = {
                column.key: getattr(user, column.key)
                for column in class_mapper(Users).columns
                # Nếu UserIDI là IDENTITY, bạn nên loại nó ra để DB tự sinh
                if column.key != 'UserIDI'
            }
            stmt = insert(generated_models.Users).values(user_dict)
            await self.session.execute(stmt)
        except IntegrityError as e:
            # Nếu xảy ra lỗi (ví dụ: Trùng Username do ràng buộc Unique)
            raise DuplicateAccountError("Tài khoản đã tồn tại trong hệ thống.")
            # (Bạn nên thay bằng Custom Exception như DuplicateAccountError)
        except Exception as e:
            raise e

    async def get_user_by_id(self, user_id: int) -> Optional[Users]:
        stmt = select(Users).where(Users.UserID == user_id)

        result = await self.session.execute(stmt)

        return result.scalars().first()
