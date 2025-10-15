from sqlalchemy import Column, Integer, String, DateTime, ForeignKey, UniqueConstraint
from sqlalchemy.sql import func
from sqlalchemy.orm import relationship
from ..core.database import Base

class ReadingProgress(Base):
    __tablename__ = "reading_progress"
    
    id = Column(Integer, primary_key=True, index=True)
    user_id = Column(Integer, ForeignKey("users.id"), nullable=False)
    book_id = Column(Integer, ForeignKey("books.id"), nullable=False)
    
    # Progress tracking
    current_chapter = Column(Integer, default=1)
    chapters_completed = Column(Integer, default=0)
    status = Column(String, default="reading")  # reading, completed, abandoned
    
    # Timestamps
    started_at = Column(DateTime(timezone=True), server_default=func.now())
    completed_at = Column(DateTime(timezone=True), nullable=True)
    last_read_at = Column(DateTime(timezone=True), server_default=func.now(), onupdate=func.now())
    
    # Relationships
    user = relationship("User")
    book = relationship("Book")
    
    __table_args__ = (
        UniqueConstraint("user_id", "book_id", name="uq_user_book_progress"),
    )
    
    def __repr__(self):
        return f"<ReadingProgress user_id={self.user_id} book_id={self.book_id} status={self.status}>"