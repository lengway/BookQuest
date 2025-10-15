from typing import List
from fastapi import APIRouter, Depends, HTTPException, status
from sqlalchemy.orm import Session

from ..core.database import get_db
from ..core.dependencies import get_current_active_user
from ..models import (
    Quiz as QuizModel,
    Question as QuestionModel,
    Option as OptionModel,
)
from ..schemas import (
    QuizCreate, QuizUpdate, Quiz as QuizSchema,
    QuestionCreate, QuestionUpdate, Question as QuestionSchema,
    OptionCreate, OptionUpdate, Option as OptionSchema,
    QuizSubmission, QuizResult, QuestionResult,
)
from ..services.grading import grade_question
from ..services.xp import compute_quiz_xp, apply_quiz_result
from datetime import datetime


router = APIRouter()


# Admin CRUD

@router.post("/", response_model=QuizSchema, status_code=status.HTTP_201_CREATED)
def create_quiz(
    quiz_in: QuizCreate,
    db: Session = Depends(get_db),
    user=Depends(get_current_active_user),
):
    if not user.is_superuser:
        raise HTTPException(status_code=status.HTTP_403_FORBIDDEN, detail="Forbidden")
    quiz = QuizModel(**quiz_in.model_dump())
    db.add(quiz)
    db.commit()
    db.refresh(quiz)
    return quiz


@router.put("/{quiz_id}", response_model=QuizSchema)
def update_quiz(
    quiz_id: int,
    quiz_in: QuizUpdate,
    db: Session = Depends(get_db),
    user=Depends(get_current_active_user),
):
    if not user.is_superuser:
        raise HTTPException(status_code=status.HTTP_403_FORBIDDEN, detail="Forbidden")
    quiz = db.query(QuizModel).filter(QuizModel.id == quiz_id).first()
    if not quiz:
        raise HTTPException(status_code=status.HTTP_404_NOT_FOUND, detail="Quiz not found")
    for field, value in quiz_in.model_dump(exclude_unset=True).items():
        setattr(quiz, field, value)
    db.add(quiz)
    db.commit()
    db.refresh(quiz)
    return quiz


@router.delete("/{quiz_id}", status_code=status.HTTP_204_NO_CONTENT)
def delete_quiz(
    quiz_id: int,
    db: Session = Depends(get_db),
    user=Depends(get_current_active_user),
):
    if not user.is_superuser:
        raise HTTPException(status_code=status.HTTP_403_FORBIDDEN, detail="Forbidden")
    quiz = db.query(QuizModel).filter(QuizModel.id == quiz_id).first()
    if not quiz:
        raise HTTPException(status_code=status.HTTP_404_NOT_FOUND, detail="Quiz not found")
    db.delete(quiz)
    db.commit()
    return None


@router.post("/{quiz_id}/questions", response_model=QuestionSchema, status_code=status.HTTP_201_CREATED)
def create_question(
    quiz_id: int,
    q_in: QuestionCreate,
    db: Session = Depends(get_db),
    user=Depends(get_current_active_user),
):
    if not user.is_superuser:
        raise HTTPException(status_code=status.HTTP_403_FORBIDDEN, detail="Forbidden")
    quiz = db.query(QuizModel).filter(QuizModel.id == quiz_id).first()
    if not quiz:
        raise HTTPException(status_code=status.HTTP_404_NOT_FOUND, detail="Quiz not found")
    q = QuestionModel(quiz_id=quiz_id, **q_in.model_dump())
    db.add(q)
    db.commit()
    db.refresh(q)
    return q


@router.put("/questions/{question_id}", response_model=QuestionSchema)
def update_question(
    question_id: int,
    q_in: QuestionUpdate,
    db: Session = Depends(get_db),
    user=Depends(get_current_active_user),
):
    if not user.is_superuser:
        raise HTTPException(status_code=status.HTTP_403_FORBIDDEN, detail="Forbidden")
    q = db.query(QuestionModel).filter(QuestionModel.id == question_id).first()
    if not q:
        raise HTTPException(status_code=status.HTTP_404_NOT_FOUND, detail="Question not found")
    for field, value in q_in.model_dump(exclude_unset=True).items():
        setattr(q, field, value)
    db.add(q)
    db.commit()
    db.refresh(q)
    return q


@router.delete("/questions/{question_id}", status_code=status.HTTP_204_NO_CONTENT)
def delete_question(
    question_id: int,
    db: Session = Depends(get_db),
    user=Depends(get_current_active_user),
):
    if not user.is_superuser:
        raise HTTPException(status_code=status.HTTP_403_FORBIDDEN, detail="Forbidden")
    q = db.query(QuestionModel).filter(QuestionModel.id == question_id).first()
    if not q:
        raise HTTPException(status_code=status.HTTP_404_NOT_FOUND, detail="Question not found")
    db.delete(q)
    db.commit()
    return None


@router.post("/questions/{question_id}/options", response_model=OptionSchema, status_code=status.HTTP_201_CREATED)
def create_option(
    question_id: int,
    o_in: OptionCreate,
    db: Session = Depends(get_db),
    user=Depends(get_current_active_user),
):
    if not user.is_superuser:
        raise HTTPException(status_code=status.HTTP_403_FORBIDDEN, detail="Forbidden")
    q = db.query(QuestionModel).filter(QuestionModel.id == question_id).first()
    if not q:
        raise HTTPException(status_code=status.HTTP_404_NOT_FOUND, detail="Question not found")
    o = OptionModel(question_id=question_id, **o_in.model_dump())
    db.add(o)
    db.commit()
    db.refresh(o)
    return o


@router.put("/options/{option_id}", response_model=OptionSchema)
def update_option(
    option_id: int,
    o_in: OptionUpdate,
    db: Session = Depends(get_db),
    user=Depends(get_current_active_user),
):
    if not user.is_superuser:
        raise HTTPException(status_code=status.HTTP_403_FORBIDDEN, detail="Forbidden")
    o = db.query(OptionModel).filter(OptionModel.id == option_id).first()
    if not o:
        raise HTTPException(status_code=status.HTTP_404_NOT_FOUND, detail="Option not found")
    for field, value in o_in.model_dump(exclude_unset=True).items():
        setattr(o, field, value)
    db.add(o)
    db.commit()
    db.refresh(o)
    return o


@router.delete("/options/{option_id}", status_code=status.HTTP_204_NO_CONTENT)
def delete_option(
    option_id: int,
    db: Session = Depends(get_db),
    user=Depends(get_current_active_user),
):
    if not user.is_superuser:
        raise HTTPException(status_code=status.HTTP_403_FORBIDDEN, detail="Forbidden")
    o = db.query(OptionModel).filter(OptionModel.id == option_id).first()
    if not o:
        raise HTTPException(status_code=status.HTTP_404_NOT_FOUND, detail="Option not found")
    db.delete(o)
    db.commit()
    return None


# Public endpoints

@router.get("/by-chapter/{chapter_id}", response_model=QuizSchema)
def get_quiz_by_chapter(chapter_id: int, db: Session = Depends(get_db), _=Depends(get_current_active_user)):
    quiz = db.query(QuizModel).filter(QuizModel.chapter_id == chapter_id, QuizModel.is_active == True).first()
    if not quiz:
        raise HTTPException(status_code=status.HTTP_404_NOT_FOUND, detail="Quiz not found")
    # скрыть ответы в options: pydantic схема уже не требует is_correct на выдачу
    return quiz


@router.post("/{quiz_id}/submit", response_model=QuizResult)
def submit_quiz(
    quiz_id: int,
    submission: QuizSubmission,
    db: Session = Depends(get_db),
    user=Depends(get_current_active_user),
):
    quiz = db.query(QuizModel).filter(QuizModel.id == quiz_id).first()
    if not quiz:
        raise HTTPException(status_code=status.HTTP_404_NOT_FOUND, detail="Quiz not found")

    # Собираем вопросы словарём
    q_map = {q.id: q for q in quiz.questions}
    results: list[QuestionResult] = []
    correct_count = 0
    total = len(quiz.questions)

    # Оценка по каждому ответу
    for qsub in submission.answers:
        q = q_map.get(qsub.question_id)
        if not q:
            continue
        payload = {}
        if qsub.type.name.lower() in ("single_choice", "multi_choice"):
            if qsub.single_multi:
                payload["selected_option_ids"] = qsub.single_multi.selected_option_ids
        elif qsub.type.name.lower() == "ordering":
            if qsub.ordering:
                payload["ordering"] = qsub.ordering.ordering
        elif qsub.type.name.lower() == "matching":
            if qsub.matching:
                payload["matches"] = qsub.matching.matches

        is_correct = grade_question(q, payload)
        results.append(QuestionResult(question_id=q.id, correct=is_correct))
        if is_correct:
            correct_count += 1

    is_perfect = total > 0 and correct_count == total
    xp_earned = compute_quiz_xp(quiz, is_perfect, user)

    # Применяем результат пользователю
    apply_quiz_result(db, user, quiz, is_perfect, xp_earned)

    # Update reading progress
    from ..models import ReadingProgress as ProgressModel
    progress = db.query(ProgressModel).filter(
        ProgressModel.user_id == user.id,
        ProgressModel.book_id == quiz.chapter.book_id
    ).first()

    if progress and is_perfect:
        # Move to next chapter if quiz passed
        if quiz.chapter.chapter_number > progress.chapters_completed:
            progress.chapters_completed = quiz.chapter.chapter_number
            progress.current_chapter = min(
                quiz.chapter.chapter_number + 1,
                quiz.chapter.book.total_chapters
            )
        
        # Check if book completed
        if progress.chapters_completed >= quiz.chapter.book.total_chapters:
            progress.status = "completed"
            progress.completed_at = datetime.utcnow()
        
        progress.last_read_at = datetime.utcnow()
        db.add(progress)

    # Сохраняем попытку
    from ..models import QuizAttempt, Answer
    attempt = QuizAttempt(
        user_id=user.id,
        quiz_id=quiz.id,
        is_perfect=is_perfect,
        score_earned=xp_earned,
        total_questions=total,
        correct_questions=correct_count,
    )
    db.add(attempt)
    db.commit()
    db.refresh(attempt)

    # Ответы
    for r in submission.answers:
        payload = {}
        if r.single_multi:
            payload["selected_option_ids"] = r.single_multi.selected_option_ids
        if r.ordering:
            payload["ordering"] = r.ordering.ordering
        if r.matching:
            payload["matches"] = r.matching.matches
        ans = Answer(
            attempt_id=attempt.id,
            question_id=r.question_id,
            selected_option_ids=payload.get("selected_option_ids"),
            ordering=payload.get("ordering"),
            matches=payload.get("matches"),
        )
        db.add(ans)

    db.commit()

    return QuizResult(
        total_questions=total,
        correct_questions=correct_count,
        is_perfect=is_perfect,
        xp_earned=xp_earned,
        question_results=results,
    )


