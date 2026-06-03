import re
from uuid import UUID

from pydantic import BaseModel, Field

from domain.schemas.exceptions import AccountBannerError, AccountDeletedError


class LoginRequest(BaseModel):
    account: str
    password: str

class LoginResponse(BaseModel):

    UserID: UUID
    Token: str
    Username: str
    EmployeeID: str
    PasswordHash: str
    UserStatus: str

    def is_login(self) -> bool:
        """0 No status,1 Active, 2 Pending, 3 Suspended, 4 Deleted"""
        match self.UserStatus:
            case 1:
                return True
            case 2:
                raise Exception("Account need validate")
            case 3:
                raise AccountBannerError()
            case 4:
                raise AccountDeletedError()

        return False;

class RegisterRequest(BaseModel):
    Username: str
    Password: str
    EmployeeID: str


    def validate_username(self) -> bool:
        return re.match(r"^(?=.{3,20}$)(?![_-])(?!.*[_-]{2})[a-zA-Z0-9_-]+(?<![_-])$", self.Username)

    def validate_password(self) -> bool:
        return re.match(r"^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[\W_])[^\s]{8,32}$", self.Password)

    def validate_account(self) -> bool:
        return self.validate_username() and self.validate_password()
