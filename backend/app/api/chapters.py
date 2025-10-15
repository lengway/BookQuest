from fastapi import APIRouter, Depends, HTTPException, status
from sqlalchemy.orm import Session

from ..core.database import get_db
from ..core.dependencies import get_current_active_user
from ..models import Chapter as ChapterModel, Book as BookModel
from ..schemas import ChapterCreate, ChapterUpdate, Chapter as ChapterSchema


router = APIRouter()


@router.post("/", response_model=ChapterSchema, status_code=status.HTTP_201_CREATED)
def create_chapter(
    chapter_in: ChapterCreate,
    db: Session = Depends(get_db),
    _user=Depends(get_current_active_user),
):
    # ensure book exists
    book = db.query(BookModel).filter(BookModel.id == chapter_in.book_id).first()
    if not book:
        raise HTTPException(status_code=status.HTTP_404_NOT_FOUND, detail="Book not found")

    chapter = ChapterModel(**chapter_in.model_dump())
    db.add(chapter)
    db.commit()
    db.refresh(chapter)
    return chapter


@router.get("/{chapter_id}", response_model=ChapterSchema)
def get_chapter(chapter_id: int, db: Session = Depends(get_db)):
    chapter = db.query(ChapterModel).filter(ChapterModel.id == chapter_id).first()
    if not chapter:
        raise HTTPException(status_code=status.HTTP_404_NOT_FOUND, detail="Chapter not found")
    return chapter


@router.put("/{chapter_id}", response_model=ChapterSchema)
def update_chapter(
    chapter_id: int,
    chapter_in: ChapterUpdate,
    db: Session = Depends(get_db),
    _user=Depends(get_current_active_user),
):
    chapter = db.query(ChapterModel).filter(ChapterModel.id == chapter_id).first()
    if not chapter:
        raise HTTPException(status_code=status.HTTP_404_NOT_FOUND, detail="Chapter not found")

    for field, value in chapter_in.model_dump(exclude_unset=True).items():
        setattr(chapter, field, value)

    db.add(chapter)
    db.commit()
    db.refresh(chapter)
    return chapter


@router.delete("/{chapter_id}", status_code=status.HTTP_204_NO_CONTENT)
def delete_chapter(
    chapter_id: int,
    db: Session = Depends(get_db),
    _user=Depends(get_current_active_user),
):
    chapter = db.query(ChapterModel).filter(ChapterModel.id == chapter_id).first()
    if not chapter:
        raise HTTPException(status_code=status.HTTP_404_NOT_FOUND, detail="Chapter not found")
    db.delete(chapter)
    db.commit()
    return None


@router.get("/", response_model=list[ChapterSchema])
def list_chapters(
    book_id: int | None = None,
    skip: int = 0,
    limit: int = 20,
    db: Session = Depends(get_db),
):
    query = db.query(ChapterModel)
    if book_id is not None:
        query = query.filter(ChapterModel.book_id == book_id)
    chapters = query.offset(skip).limit(limit).all()
    return chapters


