import uuid
from uuid import UUID

from starlette import status

from api.dependencies.auth_deps import verify_global_token
from core.container import Container
from fastapi import APIRouter, Depends, HTTPException
from dependency_injector.wiring import inject, Provide

from domain.interfaces.services.i_auth_service import IAuthService
from domain.schemas.exceptions import AccountNotFoundError, PasswordIncorrectError
from domain.schemas.service_result import ServiceResult
from domain.schemas.user_dto import LoginRequest, LoginResponse, RegisterRequest
from services import auth_service

router = APIRouter()


@router.post("/register")
@inject
async def register(request: RegisterRequest,
                   auth_service: IAuthService = Depends(Provide[Container.auth_service])):
    try:
        token = await auth_service.register(request)
        return token.to_dict()
    except Exception as e:
        raise e


@router.post("/logout")  # 1. ĐỔI THÀNH POST
@inject
async def logout(
        current_user: dict = Depends(verify_global_token),
        auth_service: IAuthService = Depends(Provide[Container.auth_service])
):
    user_id = current_user.get("sub")

    # 2. Kiểm tra cẩn thận trước khi xử lý
    if not user_id:
        raise HTTPException(
            status_code=status.HTTP_401_UNAUTHORIZED,
            detail="Token không chứa thông tin User ID."
        )

    try:
        # Ép kiểu an toàn
        user_uuid = uuid.UUID(user_id)
        result = await auth_service.logout(user_uuid)

        # Trả về kết quả
        return result.to_dict()

    except ValueError:
        # Bắt lỗi nếu user_id không phải là một UUID hợp lệ
        raise HTTPException(
            status_code=status.HTTP_400_BAD_REQUEST,
            detail="User ID không đúng định dạng UUID."
        )
    except Exception as e:
        # 3. Trả về lỗi 500 có kiểm soát (Nên có thư viện logging để ghi lại biến 'e' ở đây)
        raise HTTPException(
            status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
            detail="Đã xảy ra lỗi hệ thống khi đăng xuất."
        )
@router.post("/login")
@inject
async def login(
        request: LoginRequest,
        auth_service: IAuthService = Depends(Provide[Container.auth_service])
):
    try:
        # 1. Gọi thẳng vào hàm login của Service
        result = await auth_service.login(request)
        # 2. Nếu không có lỗi, trả về kết quả thành công
        return result.to_dict()

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
