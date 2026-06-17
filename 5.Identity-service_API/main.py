from contextlib import asynccontextmanager

from dependency_injector import containers
from fastapi import FastAPI

from api.middlewares.exception_middlewares import app_global_exception_middleware
from api.routes import router_registration
from core.container import Container


@asynccontextmanager
async def lifespan(app: FastAPI):
    container = app.container

    try:
        await container.auth_service().initADMIN()
        print("ADMIN khởi tạo thành công!")
    except Exception as e:
        print(f"Khởi tạo admin thất bại! Lỗi: {e}")

    yield
    await container.redis_service().close()
    print("🛑 Hệ thống đang tắt...")

def create_app() -> FastAPI:
    container = Container()

    app = FastAPI(
        title = "Identity Service API For OK2SHIP_2.0",
        description="Identity Service API For OK2SHIP_2.0 with fast API",
        version="26.0.0",
        contact={
            "name": "gitzno",
            "email": "hoangthuy.forjob@gmail.com",
        },
        docs_url="/docs",
        redoc_url="/redoc",
        lifespan=lifespan
    )
    app_global_exception_middleware(app)
    app.container = container

    #nhúng router tổng vào
    app.include_router(router_registration.router, prefix="/api")

    #Router health check
    @app.get("/", tags=["Root"])
    def health_check():
        return {"status": "Hello World"}

    return app

app = create_app()