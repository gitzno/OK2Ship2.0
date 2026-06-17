import uuid
from uuid import UUID

import bcrypt
from sqlalchemy import false, true
from sqlalchemy.exc import IntegrityError
from sqlalchemy.sql import cache_key

from core.UnitOfWork import UnitOfWork
from core.config import settings, Settings
from core.interfaces.cache_interface import ICacheService
from domain.interfaces.services.i_auth_service import IAuthService
from domain.interfaces.services.i_token_service import ITokenService
from domain.models.generated_models import Users
from domain.schemas.exceptions import AccountNotFoundError, PasswordIncorrectError, DuplicateAccountError, \
    ERROR_MESSAGES
from domain.schemas.service_result import ServiceResult
from domain.schemas.user_dto import LoginRequest, RegisterRequest


class AuthService(IAuthService):
    def __init__(self, uow: UnitOfWork, token_service: ITokenService, cache: ICacheService):
        self.uow = uow
        self.tokenService = token_service
        self.cache = cache



    @staticmethod
    async def __hash_password(password: str) -> str:
        password = password.encode("utf-8")

        salt = bcrypt.gensalt()

        hashed_bytes = bcrypt.hashpw(password, salt)
        return hashed_bytes.decode("utf-8")

    @staticmethod
    async def __verify_password(password: str, hashed_password: str) -> bool:
        try:
            return bcrypt.checkpw(password.encode("utf-8"), hashed_password.encode("utf-8"))
        except Exception:
            raise Exception("Hashed password error")

    async def login(self, request : LoginRequest) -> ServiceResult[str]:
        async with self.uow:
            account = request.account
            password = request.password
            user = await self.uow.users.get_user_by_account(account)

            # Kiểm tra tài khoản có tồn tại hay không
            if not user:
                raise Exception("Account not found")

            # Kiểm tra trạng thái tài khoản
            if user.UserStatus != 1:
                raise AccountNotFoundError()

            # So sánh mật khẩu
            if not await self.__verify_password(password, user.PasswordHash):
                raise PasswordIncorrectError()



            # Khởi tạo thông tin để cho vào token
            token_payload = {
                "sub": str(user.UserID),
                "username": user.Username,
                "security_stamp" : str(user.SecurityStamp),
            }



            # Khởi tạo token
            token = self.tokenService.generate_access_token(token_payload)

        return ServiceResult(True ,token, "OK")

    async def initADMIN(self) -> bool:
        async with self.uow:
            user = Users(
                UserIDI = 0,
                Username=settings.ADMIN_USERNAME,
                PasswordHash=await self.__hash_password(settings.ADMIN_PASSWORD),
                EmployeeID="ADMIN_TDMK",
                UserStatus= 1,  # pending,

            )
            try:
                try:
                    user = await self.uow.users.create_user(user)
                except DuplicateAccountError as e:
                    print("Tài khoản đã tồn tại trong hệ thống")
                try:
                    await self.uow.roles.set_role(user.UserIDI , 1)
                except DuplicateAccountError as e:
                    print("Quyền đã tồn tại trong hệ thống")
            except Exception as e:
                print(f"Lỗi khi khởi tạo {e}")
                return False
        return True

    async def register(self, request: RegisterRequest) -> ServiceResult[str]:

        # Validate data
        if not request.validate_account():
            raise AccountNotFoundError()

        user = Users(
            Username = request.Username,
            PasswordHash = await self.__hash_password(request.Password),
            EmployeeID = request.EmployeeID,
            UserStatus = 2, #pending
        )

        user = await self.user_repo.create_user(user)

        return ServiceResult(True ,user, "OK")


    async def logout(self, userId : UUID) -> ServiceResult[bool]:
        new_stamp = uuid.uuid4()
        if(userId == None):
            raise

        async with self.uow:
            try:
                is_update = await self.uow.users.update_stamp(userId, new_stamp)
            except Exception as e:
                print(e)

            if not is_update:
                return ServiceResult(False)
            await self.uow.commit()

        cache_key_auth = Settings.KEY_CACHE_AUTH_USER.format(user_ID = str(userId))

        await self.cache.set(cache_key_auth, str(new_stamp) , ttl=86400)
        return ServiceResult(True)

