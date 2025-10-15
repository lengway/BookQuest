from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware
from .core.config import settings
from .core.init_db import init_db
from .api import auth
from .api import books as books_router
from .api import chapters as chapters_router
from .api import users as users_router
from .api import quizzes as quizzes_router

# Create FastAPI app
app = FastAPI(
    title=settings.PROJECT_NAME,
    version=settings.VERSION,
    openapi_url=f"{settings.API_V1_STR}/openapi.json"
)

# Initialize database on startup
@app.on_event("startup")
async def startup_event():
    init_db()

# CORS middleware
app.add_middleware(
    CORSMiddleware,
    allow_origins=settings.BACKEND_CORS_ORIGINS,
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

@app.get("/")
async def root():
    return {
        "message": "Welcome to BookQuest API",
        "version": settings.VERSION,
        "docs": "/docs"
    }

@app.get("/health")
async def health_check():
    return {"status": "healthy"}

# Include routers
app.include_router(auth.router, prefix=f"{settings.API_V1_STR}/auth", tags=["Authentication"])
app.include_router(books_router.router, prefix=f"{settings.API_V1_STR}/books", tags=["Books"])
app.include_router(chapters_router.router, prefix=f"{settings.API_V1_STR}/chapters", tags=["Chapters"])
app.include_router(users_router.router, prefix=f"{settings.API_V1_STR}/users", tags=["Users"])
app.include_router(quizzes_router.router, prefix=f"{settings.API_V1_STR}/quizzes", tags=["Quizzes"])