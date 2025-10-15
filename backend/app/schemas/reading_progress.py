from pydantic import BaseModel
from typing import Optional
from datetime import datetime

class ReadingProgressBase(BaseModel):
    book_id: int
    current_chapter: int = 1
    chapters_completed: int = 0
    status: str = "reading"

class ReadingProgressCreate(BaseModel):
    book_id: int

class ReadingProgressUpdate(BaseModel):
    current_chapter: Optional[int] = None
    chapters_completed: Optional[int] = None
    status: Optional[str] = None

class ReadingProgress(ReadingProgressBase):
    id: int
    user_id: int
    started_at: datetime
    completed_at: Optional[datetime] = None
    last_read_at: datetime
    
    class Config:
        from_attributes = True

class ReadingProgressWithBook(ReadingProgress):
    """Progress with book details"""
    book_title: str
    book_author: str
    book_cover_url: Optional[str] = None
    total_chapters: int
    
    class Config:
        from_attributes = True