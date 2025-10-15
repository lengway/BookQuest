from sqlalchemy.orm import Session
from datetime import datetime
from ..core.config import settings
from ..models import User, Quiz

def xp_for_level(level: int) -> int:
    """Calculate XP needed to reach next level"""
    return int(500 * (level ** 1.5))

def check_level_up(user: User) -> tuple[bool, int, list[int]]:
    """
    Check if user should level up
    Returns: (leveled_up, new_level, levels_gained)
    """
    leveled_up = False
    levels_gained = []
    
    while user.current_xp >= xp_for_level(user.level):
        user.current_xp -= xp_for_level(user.level)
        user.level += 1
        levels_gained.append(user.level)
        leveled_up = True
    
    return leveled_up, user.level, levels_gained

def compute_quiz_xp(quiz: Quiz, is_perfect: bool, user: User) -> int:
    """Calculate XP earned from quiz"""
    base = quiz.quiz_xp if quiz.quiz_xp is not None else 0
    if base == 0:
        try:
            base = quiz.chapter.book.chapter_xp or 0
        except Exception:
            base = settings.QUIZ_BASE_XP
    
    bonus = 0
    if is_perfect and settings.STREAK_STEP > 0:
        next_streak = (user.reading_streak or 0) + 1
        if next_streak % settings.STREAK_STEP == 0:
            bonus = settings.STREAK_BONUS_XP
    
    return base + bonus

def apply_quiz_result(db: Session, user: User, quiz: Quiz, is_perfect: bool, xp_earned: int) -> dict:
    """
    Apply quiz results to user (XP, streak, level)
    Returns dict with level_up info
    """
    # Update streak
    if is_perfect:
        user.reading_streak = (user.reading_streak or 0) + 1
    else:
        user.reading_streak = 0
    
    # Add XP
    user.current_xp = (user.current_xp or 0) + xp_earned
    user.total_xp = (user.total_xp or 0) + xp_earned
    
    # Check for level up
    leveled_up, new_level, levels_gained = check_level_up(user)
    
    # Update last reading date
    user.last_reading_date = datetime.utcnow()
    
    db.add(user)
    
    return {
        "leveled_up": leveled_up,
        "new_level": new_level,
        "levels_gained": levels_gained,
    }