from fastapi import APIRouter, Depends, HTTPException, status
from sqlalchemy.orm import Session

from ..core.database import get_db
from ..core.dependencies import get_current_active_user
from ..core.security import get_password_hash
from ..models import User as UserModel
from ..schemas import UserResponse, UserUpdate


router = APIRouter()


@router.get("/", response_model=list[UserResponse])
def list_users(skip: int = 0, limit: int = 20, db: Session = Depends(get_db), _=Depends(get_current_active_user)):
    users = db.query(UserModel).offset(skip).limit(limit).all()
    return users


@router.get("/{user_id}", response_model=UserResponse)
def get_user(user_id: int, db: Session = Depends(get_db), _=Depends(get_current_active_user)):
    user = db.query(UserModel).filter(UserModel.id == user_id).first()
    if not user:
        raise HTTPException(status_code=status.HTTP_404_NOT_FOUND, detail="User not found")
    return user


@router.put("/{user_id}", response_model=UserResponse)
def update_user(
    user_id: int,
    user_in: UserUpdate,
    db: Session = Depends(get_db),
    _=Depends(get_current_active_user),
):
    user = db.query(UserModel).filter(UserModel.id == user_id).first()
    if not user:
        raise HTTPException(status_code=status.HTTP_404_NOT_FOUND, detail="User not found")

    update_data = user_in.model_dump(exclude_unset=True)

    # uniqueness checks for email and username
    if "email" in update_data:
        exists = db.query(UserModel).filter(UserModel.email == update_data["email"], UserModel.id != user_id).first()
        if exists:
            raise HTTPException(status_code=status.HTTP_400_BAD_REQUEST, detail="Email already in use")

    if "username" in update_data:
        exists = db.query(UserModel).filter(UserModel.username == update_data["username"], UserModel.id != user_id).first()
        if exists:
            raise HTTPException(status_code=status.HTTP_400_BAD_REQUEST, detail="Username already taken")

    # hash password if provided
    if "password" in update_data and update_data["password"]:
        update_data.pop("password")
        user.hashed_password = get_password_hash(user_in.password)  # type: ignore[arg-type]

    for field, value in update_data.items():
        setattr(user, field, value)

    db.add(user)
    db.commit()
    db.refresh(user)
    return user


@router.delete("/{user_id}", status_code=status.HTTP_204_NO_CONTENT)
def delete_user(user_id: int, db: Session = Depends(get_db), _=Depends(get_current_active_user)):
    user = db.query(UserModel).filter(UserModel.id == user_id).first()
    if not user:
        raise HTTPException(status_code=status.HTTP_404_NOT_FOUND, detail="User not found")
    db.delete(user)
    db.commit()
    return None


