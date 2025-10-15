from typing import Optional
from pydantic import BaseModel, Field


class ChapterBase(BaseModel):
    book_id: int
    chapter_number: int = Field(..., ge=1)
    title: str = Field(..., min_length=1)
    content: str = Field(..., min_length=1)
    estimated_reading_time: Optional[int] = None  # minutes


class ChapterCreate(ChapterBase):
    pass


class ChapterUpdate(BaseModel):
    chapter_number: Optional[int] = Field(default=None, ge=1)
    title: Optional[str] = None
    content: Optional[str] = None
    estimated_reading_time: Optional[int] = None


class Chapter(ChapterBase):
    id: int

    class Config:
        from_attributes = True


