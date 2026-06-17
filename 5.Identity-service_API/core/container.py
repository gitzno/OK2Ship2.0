from dependency_injector import containers, providers
from redis.cache import CacheFactory

from core.UnitOfWork import UnitOfWork
from core.config import settings
from core.database import async_session
from infrastructures.cache import redis_service
from infrastructures.cache.redis_service import RedisService
from services.auth_service import AuthService
from services.token_service import TokenService


class Container(containers.DeclarativeContainer):


    wiring_config = containers.WiringConfiguration(packages=[
        "api.routes.v1",
        "api.dependencies"
    ])


    """
    providers.Factory: Mỗi lần bạn gọi, nó sẽ tạo ra một đối tượng mới tinh. (Rất hợp cho Repository và Service).

    providers.Singleton: Tạo ra duy nhất một đối tượng và dùng chung cho toàn bộ app. (Hợp cho các class cấu hình hệ thống).

    providers.Resource: Chuyên dùng để quản lý các tài nguyên cần thao tác Đóng/Mở (như kết nối Database, File, Network). 
    """
    session_factory = providers.Singleton(lambda: async_session)

    #Register Repository
    uow = providers.Factory(
        UnitOfWork,
        session_factory=session_factory
    )

    #Register Service
    token_service = providers.Singleton(TokenService)
    redis_service = providers.Singleton(RedisService, redis_url = settings.REDIS_URL)

    auth_service = providers.Factory(
        AuthService,
        uow=uow,
        token_service=token_service,
        cache=redis_service
    )


