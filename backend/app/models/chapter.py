from sqlalchemy import Column, Integer, String, Text, ForeignKey
from sqlalchemy.orm import relationship
from ..core.database import Base

class Chapter(Base):
    __tablename__ = "chapters"

    id = Column(Integer, primary_key=True, index=True)
    book_id = Column(Integer, ForeignKey("books.id"), nullable=False)
    chapter_number = Column(Integer, nullable=False)
    title = Column(String, nullable=False)
    content = Column(Text, nullable=False)
    estimated_reading_time = Column(Integer)  # в минутах
    
    # Relationships
    book = relationship("Book", back_populates="chapters")
    quiz = relationship("Quiz", back_populates="chapter", uselist=False)

    def __repr__(self):
        return f"<Chapter {self.chapter_number}: {self.title}>" 