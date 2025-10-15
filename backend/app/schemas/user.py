from pydantic import BaseModel, EmailStr, Field
from typing import Optional
from datetime import datetime

# Base schema
class UserBase(BaseModel):
    email: EmailStr
    username: str = Field(..., min_length=3, max_length=50)

# Schema for user registration
class UserCreate(UserBase):
    password: str = Field(..., min_length=6, max_length=100)

# Schema for user login
class UserLogin(BaseModel):
    email: EmailStr
    password: str

# Schema for user update
class UserUpdate(BaseModel):
    email: Optional[EmailStr] = None
    username: Optional[str] = Field(default=None, min_length=3, max_length=50)
    password: Optional[str] = Field(default=None, min_length=6, max_length=100)
    is_active: Optional[bool] = None

# Schema for response (what we return to client)
class UserResponse(UserBase):
    id: int
    level: int
    current_xp: int
    total_xp: int
    reading_streak: int
    created_at: datetime
    is_active: bool
    
    class Config:
        from_attributes = True

# Schema for user profile/stats
class UserStats(BaseModel):
    id: int
    username: str
    level: int
    current_xp: int
    total_xp: int
    reading_streak: int
    books_read: int = 0  # Will calculate later
    badges_count: int = 0  # Will calculate later
    
    class Config:
        from_attributes = True

# Token schemas
class Token(BaseModel):
    access_token: str
    token_type: str = "bearer"

class TokenData(BaseModel):
    user_id: Optional[int] = None