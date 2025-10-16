import { BrowserRouter, Routes, Route } from "react-router-dom";
import { AuthProvider } from "./auth/AuthContext";
import ProtectedRoute from "./auth/ProtectedRoute";
import Login from "./pages/Login";
import BooksList from "./pages/Books/List";
import BookCreate from "./pages/Books/Create";
import BookEdit from "./pages/Books/Edit";
import ChaptersList from "./pages/Chapters/List";
import ChapterCreate from "./pages/Chapters/Create";
import QuizBuilder from "./pages/Quizzes/QuizBuilder";
import QuizzesList from "./pages/Quizzes/List";
import QuestionEdit from "./pages/Quizzes/QuestionEdit";
import Home from "./pages/Home";
import ChapterEdit from "./pages/Chapters/Edit";

export default function App() {
  return (
    <AuthProvider>
      <BrowserRouter>
        <Routes>
          <Route path="/login" element={<Login />} />
          <Route
            path="/"
            element={
              <ProtectedRoute>
                <Home />
              </ProtectedRoute>
            }
          />
          <Route
            path="/books"
            element={
              <ProtectedRoute>
                <BooksList />
              </ProtectedRoute>
            }
          />
          <Route
            path="/books/new"
            element={
              <ProtectedRoute>
                <BookCreate />
              </ProtectedRoute>
            }
          />
          <Route
            path="/books/:id"
            element={
              <ProtectedRoute>
                <BookEdit />
              </ProtectedRoute>
            }
          />
          <Route
            path="/chapters"
            element={
              <ProtectedRoute>
                <ChaptersList />
              </ProtectedRoute>
            }
          />
          <Route
            path="/chapters/new"
            element={
              <ProtectedRoute>
                <ChapterCreate />
              </ProtectedRoute>
            }
          />
          <Route
            path="/quizzes"
            element={
              <ProtectedRoute>
                <QuizzesList />
              </ProtectedRoute>
            }
          />
          <Route
            path="/quizzes/:chapterId/questions/:questionId"
            element={
              <ProtectedRoute>
                <QuestionEdit />
              </ProtectedRoute>
            }
          />
          <Route
            path="/quizzes/:chapterId"
            element={
              <ProtectedRoute>
                <QuizBuilder />
              </ProtectedRoute>
            }
          />
          <Route
            path="/chapters/:id"
            element={
              <ProtectedRoute>
                <ChapterEdit />
              </ProtectedRoute>
            }
          />
        </Routes>
      </BrowserRouter>
    </AuthProvider>
  );
}