from sqlalchemy import Column, Integer, String, Text, ForeignKey, Boolean, DateTime, Enum as SQLEnum, UniqueConstraint
from sqlalchemy.dialects.postgresql import JSONB
from sqlalchemy.orm import relationship
from sqlalchemy.sql import func
import enum
from ..core.database import Base


class QuestionType(str, enum.Enum):
    SINGLE_CHOICE = "single_choice"
    MULTI_CHOICE = "multi_choice"
    ORDERING = "ordering"
    MATCHING = "matching"


class Quiz(Base):
    __tablename__ = "quizzes"

    id = Column(Integer, primary_key=True, index=True)
    chapter_id = Column(Integer, ForeignKey("chapters.id"), nullable=False, unique=True)
    title = Column(String, nullable=False)
    description = Column(Text)
    is_active = Column(Boolean, default=True)
    quiz_xp = Column(Integer, nullable=True)

    created_at = Column(DateTime(timezone=True), server_default=func.now())
    updated_at = Column(DateTime(timezone=True), onupdate=func.now())

    # relationships
    chapter = relationship("Chapter", back_populates="quiz")
    questions = relationship("Question", back_populates="quiz", cascade="all, delete-orphan")


class Question(Base):
    __tablename__ = "questions"

    id = Column(Integer, primary_key=True, index=True)
    quiz_id = Column(Integer, ForeignKey("quizzes.id"), nullable=False)
    type = Column(SQLEnum(QuestionType), nullable=False)
    text = Column(Text, nullable=False)
    order_index = Column(Integer, default=0)
    score = Column(Integer, default=1)

    quiz = relationship("Quiz", back_populates="questions")
    options = relationship("Option", back_populates="question", cascade="all, delete-orphan")


class Option(Base):
    __tablename__ = "options"

    id = Column(Integer, primary_key=True, index=True)
    question_id = Column(Integer, ForeignKey("questions.id"), nullable=False)
    text = Column(Text, nullable=False)
    is_correct = Column(Boolean, default=False)
    order_index = Column(Integer, nullable=True)
    match_key = Column(String, nullable=True)

    question = relationship("Question", back_populates="options")

    __table_args__ = (
        # для ordering можно потребовать уникальный порядок в рамках вопроса (не строго обязательно)
        UniqueConstraint("question_id", "order_index", name="uq_option_order_per_question"),
    )


class QuizAttempt(Base):
    __tablename__ = "quiz_attempts"

    id = Column(Integer, primary_key=True, index=True)
    user_id = Column(Integer, ForeignKey("users.id"), nullable=False)
    quiz_id = Column(Integer, ForeignKey("quizzes.id"), nullable=False)
    started_at = Column(DateTime(timezone=True), server_default=func.now())
    finished_at = Column(DateTime(timezone=True))
    is_perfect = Column(Boolean, default=False)
    score_earned = Column(Integer, default=0)
    total_questions = Column(Integer, default=0)
    correct_questions = Column(Integer, default=0)

    # relationships
    answers = relationship("Answer", back_populates="attempt", cascade="all, delete-orphan")
    quiz = relationship("Quiz")


class Answer(Base):
    __tablename__ = "answers"

    id = Column(Integer, primary_key=True, index=True)
    attempt_id = Column(Integer, ForeignKey("quiz_attempts.id"), nullable=False)
    question_id = Column(Integer, ForeignKey("questions.id"), nullable=False)
    # хранение ответа в формате JSON, в зависимости от типа вопроса
    selected_option_ids = Column(JSONB, nullable=True)  # list[int]
    ordering = Column(JSONB, nullable=True)             # list[int] (option ids в правильной последовательности)
    matches = Column(JSONB, nullable=True)              # dict[str, str] или dict[option_id, option_id]

    attempt = relationship("QuizAttempt", back_populates="answers")

