import bcrypt

from domain.interfaces.repositories.i_user_repository import IUserRepository
from domain.interfaces.services.i_auth_service import IAuthService
from domain.interfaces.services.i_token_service import ITokenService
from domain.models.generated_models import Users
from domain.schemas.exceptions import AccountNotFoundError, PasswordIncorrectError
from domain.schemas.user_dto import LoginRequest, RegisterRequest


class AuthService(IAuthService):
    def __init__(self, user_repo: IUserRepository, token_service: ITokenService):
        self.user_repo = user_repo
        self.tokenService = token_service

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

    async def login(self, request : LoginRequest) -> str:

        account = request.account
        password = request.password
        user = await self.user_repo.get_user_by_account(account)

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
            "sub": user.UserID,
            "username": user.Username,
        }

        # Khởi tạo token
        token = self.tokenService.generate_access_token(token_payload)

        return token

    async def register(self, request: RegisterRequest) -> Users:

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

        return user



