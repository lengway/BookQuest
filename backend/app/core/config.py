from pydantic_settings import BaseSettings
from typing import Optional

class Settings(BaseSettings):
    # App
    PROJECT_NAME: str = "BookQuest API"
    VERSION: str = "1.0.0"
    API_V1_STR: str = "/api/v1"
    
    # Database
    POSTGRES_USER: str = "bookquest"
    POSTGRES_PASSWORD: str = "bookquest123"
    POSTGRES_HOST: str = "localhost"
    POSTGRES_PORT: str = "5432"
    POSTGRES_DB: str = "bookquest_db"
    
    @property
    def DATABASE_URL(self) -> str:
        return f"postgresql://{self.POSTGRES_USER}:{self.POSTGRES_PASSWORD}@{self.POSTGRES_HOST}:{self.POSTGRES_PORT}/{self.POSTGRES_DB}"
    
    # Security
    SECRET_KEY: str = "your-secret-key-change-this-in-production-make-it-long-and-random"
    ALGORITHM: str = "HS256"
    ACCESS_TOKEN_EXPIRE_MINUTES: int = 60 * 24 * 7  # 7 days
    
    # CORS
    BACKEND_CORS_ORIGINS: list = ["*"]  # В продакшене поменяй на конкретные домены

    # Quiz/XP settings
    QUIZ_BASE_XP: int = 100
    STREAK_STEP: int = 3
    STREAK_BONUS_XP: int = 50
    
    class Config:
        case_sensitive = True
        env_file = ".env"

settings = Settings()