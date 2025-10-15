from typing import Optional
from datetime import datetime
from pydantic import BaseModel, Field

from ..models.book import DifficultyLevel

# Base schema reflecting the Book model fields
class BookBase(BaseModel):
    title: str = Field(..., min_length=1)
    author: str = Field(..., min_length=1)
    description: Optional[str] = None
    cover_image_url: Optional[str] = None

    # Metadata
    genre: Optional[str] = None
    difficulty: DifficultyLevel = DifficultyLevel.INTERMEDIATE
    total_chapters: int = 0
    estimated_reading_time: Optional[int] = None  # minutes
    language: str = "ru"

    # XP rewards
    chapter_xp: int = 100
    completion_xp: int = 500


class BookCreate(BookBase):
    pass


class BookUpdate(BaseModel):
    title: Optional[str] = None
    author: Optional[str] = None
    description: Optional[str] = None
    cover_image_url: Optional[str] = None
    genre: Optional[str] = None
    difficulty: Optional[DifficultyLevel] = None
    total_chapters: Optional[int] = None
    estimated_reading_time: Optional[int] = None
    language: Optional[str] = None
    chapter_xp: Optional[int] = None
    completion_xp: Optional[int] = None


class BookInDBBase(BookBase):
    id: int
    created_at: Optional[datetime] = None
    updated_at: Optional[datetime] = None

    class Config:
        from_attributes = True


class Book(BookInDBBase):
    pass