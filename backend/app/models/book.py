from sqlalchemy import Column, Integer, String, Text, DateTime, Enum as SQLEnum
from sqlalchemy.sql import func
from sqlalchemy.orm import relationship
import enum
from ..core.database import Base

class DifficultyLevel(str, enum.Enum):
    BEGINNER = "beginner"
    INTERMEDIATE = "intermediate"
    ADVANCED = "advanced"

class Book(Base):
    __tablename__ = "books"

    id = Column(Integer, primary_key=True, index=True)
    title = Column(String, nullable=False, index=True)
    author = Column(String, nullable=False)
    description = Column(Text)
    cover_image_url = Column(String)
    
    # Metadata
    genre = Column(String)
    difficulty = Column(SQLEnum(DifficultyLevel), default=DifficultyLevel.INTERMEDIATE)
    total_chapters = Column(Integer, default=0)
    estimated_reading_time = Column(Integer)  # в минутах
    language = Column(String, default="ru")  # ru, kk, en
    
    # XP rewards
    chapter_xp = Column(Integer, default=100)  # XP за главу
    completion_xp = Column(Integer, default=500)  # XP за всю книгу
    
    # Timestamps
    created_at = Column(DateTime(timezone=True), server_default=func.now())
    updated_at = Column(DateTime(timezone=True), onupdate=func.now())
    
    # Relationships
    chapters = relationship("Chapter", back_populates="book", cascade="all, delete-orphan")

    def __repr__(self):
        return f"<Book {self.title} by {self.author}>"