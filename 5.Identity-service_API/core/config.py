from typing import ClassVar

from pydantic import computed_field
from pydantic_settings import BaseSettings, SettingsConfigDict

class Settings(BaseSettings):
    PROJECT_NAME: str = "My Project"

    #Cấu hình DB
    DB_USER: str
    DB_PASSWORD: str
    DB_SERVER: str
    DB_PORT: int = 1433
    DB_NAME: str
    DB_DRIVER: str = "ODBC Driver 17 for SQL Server"

    KEY_CACHE_AUTH_USER: ClassVar[str] = "user_{user_ID}"


    @computed_field
    @property
    def DATA_TABLE_URL(self) -> str:
        driver_parsed = self.DB_DRIVER.upper().replace(" ", "+")

        return(
            f"mssql+aioodbc://{self.DB_USER}:{self.DB_PASSWORD}@"
            f"{self.DB_SERVER}:{self.DB_PORT}/{self.DB_NAME}?"
            f"driver={driver_parsed}&TrustServerCertificate=yes"
        )

    model_config = SettingsConfigDict(env_file=".env", env_file_encoding="utf-8")

    # Cấu hình bảo mật JWT
    JWT_SECRET_KEY: str
    JWT_ALGORITHM: str = "HS256"
    JWT_ACCESS_TOKEN_EXPIRE_MINUTES: int = 120

    ADMIN_USERNAME: str = "ADMIN_USERNAME"
    ADMIN_PASSWORD: str = "ADMIN_PASSWORD"
    REDIS_URL_PORT: str = "REDIS_URL_PORT"
    REDIS_HOST: str = "localhost"
    REDIS_DB: int = 0
    @computed_field
    @property
    def REDIS_URL(self) -> str:
        return f"redis://{self.REDIS_HOST}:{self.REDIS_URL_PORT}/{self.REDIS_DB}"

settings = Settings()