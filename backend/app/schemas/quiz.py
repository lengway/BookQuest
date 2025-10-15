from typing import List, Optional, Dict
from pydantic import BaseModel, Field
from ..models.quiz import QuestionType


class OptionBase(BaseModel):
    text: str
    is_correct: Optional[bool] = None  # не выдаём наружу в публичной выдаче
    order_index: Optional[int] = None
    match_key: Optional[str] = None


class OptionCreate(OptionBase):
    pass


class OptionUpdate(BaseModel):
    text: Optional[str] = None
    is_correct: Optional[bool] = None
    order_index: Optional[int] = None
    match_key: Optional[str] = None


class Option(OptionBase):
    id: int

    class Config:
        from_attributes = True


class QuestionBase(BaseModel):
    type: QuestionType
    text: str
    order_index: int = 0
    score: int = 1


class QuestionCreate(QuestionBase):
    pass


class QuestionUpdate(BaseModel):
    type: Optional[QuestionType] = None
    text: Optional[str] = None
    order_index: Optional[int] = None
    score: Optional[int] = None


class Question(QuestionBase):
    id: int
    options: List[Option] = []

    class Config:
        from_attributes = True


class QuizBase(BaseModel):
    chapter_id: int
    title: str
    description: Optional[str] = None
    is_active: bool = True
    quiz_xp: Optional[int] = None


class QuizCreate(QuizBase):
    pass


class QuizUpdate(BaseModel):
    title: Optional[str] = None
    description: Optional[str] = None
    is_active: Optional[bool] = None
    quiz_xp: Optional[int] = None


class Quiz(QuizBase):
    id: int
    questions: List[Question] = []

    class Config:
        from_attributes = True


# Submission models

class SingleMultiAnswer(BaseModel):
    selected_option_ids: List[int]


class OrderingAnswer(BaseModel):
    ordering: List[int]  # option ids в порядке


class MatchingAnswer(BaseModel):
    matches: Dict[str, str]  # можно также int->int, выберем строки для совместимости


class QuestionSubmission(BaseModel):
    question_id: int
    type: QuestionType
    single_multi: Optional[SingleMultiAnswer] = None
    ordering: Optional[OrderingAnswer] = None
    matching: Optional[MatchingAnswer] = None


class QuizSubmission(BaseModel):
    answers: List[QuestionSubmission]


class QuestionResult(BaseModel):
    question_id: int
    correct: bool


class QuizResult(BaseModel):
    total_questions: int
    correct_questions: int
    is_perfect: bool
    xp_earned: int
    question_results: List[QuestionResult]


