from sqlalchemy import inspect
from .database import Base, engine

def init_db():
    """Initialize database tables"""
    # Import all models here to register them with Base
    from ..models import User, Book, Chapter, Quiz, Question, Option, QuizAttempt, Answer
    
    # Create all tables
    Base.metadata.create_all(bind=engine)
    
    # Check what tables were created
    inspector = inspect(engine)
    tables = inspector.get_table_names()
    print(f"✅ Database initialized. Tables: {tables}")
    
    return tables