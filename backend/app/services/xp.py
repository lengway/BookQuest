from sqlalchemy.orm import Session
from datetime import datetime
from ..core.config import settings
from ..models import User, Quiz


def compute_quiz_xp(quiz: Quiz, is_perfect: bool, user: User) -> int:
    base = quiz.quiz_xp if quiz.quiz_xp is not None else 0
    if base == 0:
        # fallback: используем XP главы книги, если доступно
        # quiz.chapter.book.chapter_xp может не быть загружен; пусть будет 0 при отсутствии
        try:
            base = quiz.chapter.book.chapter_xp or 0
        except Exception:
            base = settings.QUIZ_BASE_XP

    bonus = 0
    if is_perfect and settings.STREAK_STEP > 0:
        # Бонус за каждый кратный STREAK_STEP стрик
        next_streak = (user.reading_streak or 0) + 1
        if next_streak % settings.STREAK_STEP == 0:
            bonus = settings.STREAK_BONUS_XP
    return base + bonus


def apply_quiz_result(db: Session, user: User, quiz: Quiz, is_perfect: bool, xp_earned: int) -> None:
    # streak update
    if is_perfect:
        user.reading_streak = (user.reading_streak or 0) + 1
    else:
        user.reading_streak = 0

    # xp update
    user.current_xp = (user.current_xp or 0) + xp_earned
    user.total_xp = (user.total_xp or 0) + xp_earned

    db.add(user)


