from starlette import status

from core.container import Container
from fastapi import APIRouter, Depends, HTTPException
from dependency_injector.wiring import inject, Provide

from domain.interfaces.services.i_auth_service import IAuthService
from domain.schemas.exceptions import AccountNotFoundError, PasswordIncorrectError
from domain.schemas.user_dto import LoginRequest, LoginResponse, RegisterRequest

router = APIRouter()


@router.post("/register", response_model=RegisterRequest, tags=["auth"])
@inject
async def register(request: RegisterRequest,
                   auth_service: IAuthService = Depends(Provide[Container.auth_service])):
    try:
        token = await auth_service.register(request)
        return token
    except Exception as e:
        raise e

@router.post("/login", response_model=LoginResponse, tags=["auth"])
@inject
async def login(
        request: LoginRequest,
        auth_service: IAuthService = Depends(Provide[Container.auth_service])
):
    try:
        # 1. Gọi thẳng vào hàm login của Service
        token = await auth_service.login(request)

        # 2. Nếu không có lỗi, trả về kết quả thành công
        return LoginResponse(
            access_token=token,
            token_type="bearer",
            message="Đăng nhập thành công!"
        )

    # 3. Bắt các lỗi Nghiệp vụ (Domain Exceptions) và dịch nó thành lỗi HTTP
    except AccountNotFoundError:
        raise HTTPException(
            status_code=status.HTTP_403_FORBIDDEN,
            detail="Tài khoản chưa được kích hoạt hoặc đã bị khóa."
        )

    except PasswordIncorrectError:
        raise HTTPException(
            status_code=status.HTTP_401_UNAUTHORIZED,
            detail="Tài khoản hoặc mật khẩu không chính xác."
        )

    except Exception as e:
        # Bắt trường hợp "Account not found" mà bạn ném ra dưới dạng Exception thường
        if str(e) == "Account not found":
            raise HTTPException(
                status_code=status.HTTP_404_NOT_FOUND,
                detail="Tài khoản không tồn tại."
            )

        # Lỗi hệ thống ngoài ý muốn (DB sập, sai logic thuật toán...)
        raise HTTPException(
            status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
            detail=f"Đã xảy ra lỗi hệ thống, vui lòng thử lại sau. {repr(e)}"
        )
