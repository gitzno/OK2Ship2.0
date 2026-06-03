class DomainException(Exception):
    def __init__(self, error_code: str, custom_message: str = None):
        self.error_code = error_code

        error_info = ERROR_MESSAGES[error_code]

        self.user_message = error_info["user_message"]

        # Dev message gốc + nối thêm thông tin chi tiết động từ nơi ném lỗi (nếu có)
        self.dev_message = error_info["dev_message"]
        if custom_message:
            self.dev_message += f" Details: {custom_message}"
        super().__init__(self.dev_message)


ERROR_MESSAGES = {
    "ACCOUNT_BANNED": {
        "dev_message": "This account was banned",
        "user_message": "Tài khoản của bạn đã bị cấm!"
    },

    "ACCOUNT_DELETED": {
        "dev_message": "This account was deleted",
        "user_message": "Tài khoản của bạn không tồn tại!"
    },

    "ACCOUNT_NOT_FOUND": {
        "dev_message": "This account was not found",
        "user_message": "Tài khoản của bạn không tồn tại!"
    },

    "PASSWORD_INCORRECT": {
        "dev_message": "This account was not found",
        "user_message": "Tài khoản của bạn không tồn tại!"
    },

    "TOKEN_EXPIRED": {
        "dev_message": "This session has expired",
        "user_message": "Phiên đăng nhập hết hạn xin hãy đăng nhập lại!"
    },

    "TOKEN_INVALID": {
        "dev_message": "This token is invalid, Please confirm again",
        "user_message": "Phiên đăng nhập hết hạn xin hãy đăng nhập lại!"
    },

    "USER_DUPLICATE": {
        "dev_message": "This account is valid, Please confirm again",
        "user_message": "Tài khoản này đã tồn tại"
    },

    "ACCOUNT_VALIDATION": {
        "dev_message": "Account not match with validation",
        "user_message": "Tài khoản mật khẩu không phù hợp."
    }
}


class AccountBannerError(DomainException):
    def __init__(self, custom_message: str = None):
        super().__init__(error_code="ACCOUNT_BANNED", custom_message=custom_message)

class AccountDeletedError(DomainException):
    def __init__(self, custom_message: str = None):
        super().__init__(error_code="ACCOUNT_DELETED", custom_message=custom_message)


class AccountNotFoundError(DomainException):
    def __init__(self, custom_message: str = None):
        super().__init__(error_code="ACCOUNT_NOT_FOUND", custom_message=custom_message)

class PasswordIncorrectError(DomainException):
    def __init__(self, custom_message: str = None):
        super().__init__(error_code="PASSWORD_INCORRECT", custom_message=custom_message)

class TokenExpiredError(DomainException):
    def __init__(self, custom_message: str = None):
        super().__init__(error_code="TOKEN_EXPIRED", custom_message=custom_message)

class TokenInvalidError(DomainException):
    def __init__(self, custom_message: str = None):
        super().__init__(error_code="TOKEN_INVALID", custom_message=custom_message)

class DuplicateAccountError(DomainException):
    def __init__(self, custom_message: str = None):
        super().__init__(error_code="USER_DUPLICATE", custom_message=custom_message)