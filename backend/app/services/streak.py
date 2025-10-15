from datetime import datetime
from ..models import User

def update_reading_streak(user: User) -> dict:
    """
    Update user's reading streak based on last_reading_date
    Returns dict with streak info
    """
    today = datetime.utcnow().date()
    
    if user.last_reading_date is None:
        # First time reading
        user.reading_streak = 1
        user.last_reading_date = datetime.utcnow()
        return {
            "streak": 1,
            "streak_continued": True,
            "first_read": True
        }
    
    last_read = user.last_reading_date.date()
    days_diff = (today - last_read).days
    
    if days_diff == 0:
        # Already read today, no change
        return {
            "streak": user.reading_streak,
            "streak_continued": False,
            "already_read_today": True
        }
    elif days_diff == 1:
        # Continue streak
        user.reading_streak += 1
        user.last_reading_date = datetime.utcnow()
        return {
            "streak": user.reading_streak,
            "streak_continued": True,
            "streak_broken": False
        }
    else:
        # Streak broken
        old_streak = user.reading_streak
        user.reading_streak = 1
        user.last_reading_date = datetime.utcnow()
        return {
            "streak": 1,
            "streak_continued": False,
            "streak_broken": True,
            "old_streak": old_streak
        }

def get_streak_status(user: User) -> dict:
    """Get current streak status without updating"""
    if user.last_reading_date is None:
        return {
            "current_streak": 0,
            "at_risk": False,
            "broken": False
        }
    
    today = datetime.utcnow().date()
    last_read = user.last_reading_date.date()
    days_diff = (today - last_read).days
    
    return {
        "current_streak": user.reading_streak,
        "at_risk": days_diff == 1,  # Didn't read today yet
        "broken": days_diff > 1,
        "days_since_last_read": days_diff
    }