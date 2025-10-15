from fastapi import APIRouter, Depends, HTTPException, status
from sqlalchemy.orm import Session

from ..core.database import get_db
from ..core.dependencies import get_current_active_user
from ..models import Book as BookModel
from ..schemas import BookCreate, BookUpdate, Book as BookSchema


router = APIRouter()


@router.post("/", response_model=BookSchema, status_code=status.HTTP_201_CREATED)
def create_book(
    book_in: BookCreate,
    db: Session = Depends(get_db),
    _user=Depends(get_current_active_user),
):
    book = BookModel(**book_in.model_dump())
    db.add(book)
    db.commit()
    db.refresh(book)
    return book


@router.get("/{book_id}", response_model=BookSchema)
def get_book(book_id: int, db: Session = Depends(get_db)):
    book = db.query(BookModel).filter(BookModel.id == book_id).first()
    if not book:
        raise HTTPException(status_code=status.HTTP_404_NOT_FOUND, detail="Book not found")
    return book


@router.put("/{book_id}", response_model=BookSchema)
def update_book(
    book_id: int,
    book_in: BookUpdate,
    db: Session = Depends(get_db),
    _user=Depends(get_current_active_user),
):
    book = db.query(BookModel).filter(BookModel.id == book_id).first()
    if not book:
        raise HTTPException(status_code=status.HTTP_404_NOT_FOUND, detail="Book not found")

    for field, value in book_in.model_dump(exclude_unset=True).items():
        setattr(book, field, value)

    db.add(book)
    db.commit()
    db.refresh(book)
    return book


@router.delete("/{book_id}", status_code=status.HTTP_204_NO_CONTENT)
def delete_book(
    book_id: int,
    db: Session = Depends(get_db),
    _user=Depends(get_current_active_user),
):
    book = db.query(BookModel).filter(BookModel.id == book_id).first()
    if not book:
        raise HTTPException(status_code=status.HTTP_404_NOT_FOUND, detail="Book not found")
    db.delete(book)
    db.commit()
    return None


@router.get("/", response_model=list[BookSchema])
def list_books(skip: int = 0, limit: int = 20, db: Session = Depends(get_db)):
    books = db.query(BookModel).offset(skip).limit(limit).all()
    return books


