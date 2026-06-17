from abc import abstractmethod, ABC
from typing import Optional
from uuid import UUID

from sqlalchemy.ext.mypy.util import fail

from domain.models.generated_models import Users
from domain.schemas.service_result import ServiceResult
from domain.schemas.user_dto import LoginResponse


class IUserRepository(ABC):



    @abstractmethod
    async def get_user_by_account(self, account: str) -> Optional[Users]:
        """
        Tìm kiếm thông tin người dùng dựa trên tên tài khoản.

        Args:
            account (str): Tên đăng nhập hoặc email của người dùng.

        Returns:
            Optional[UserModel]: Trả về đối tượng Model người dùng nếu tìm thấy, ngược lại trả về None.
        """
        pass

    @abstractmethod
    async def create_user(self, user: Users) -> Users:
        """
        Tạo người dùng bằng thông tin của user

        Args:
            user (User): thông tin của người dùng.

        Returns:
            Optional[Users]: Trả về thông tin người dùng vừa tạo
        """

    @abstractmethod
    async def get_by_id(self, user_id: UUID) -> Optional[Users]:
        """
        Lấy thông tin đăng nhập người dùng bằng UUID định danh
        :param user_id: user_id người dùng

        :return:
        User thông tin đăng nhập của người dùng
        """
        pass

    @abstractmethod
    async def update_stamp(self,  user_id: UUID,stamp: UUID) -> bool:
        """
        Cập nhật stamp mới cho user
        :param stamp: stamp
        :param user_id: id của người dùng
        :return:trạng thái của việc update
        """
        pass
