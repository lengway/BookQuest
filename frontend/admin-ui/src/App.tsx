import { BrowserRouter, Routes, Route, Navigate } from "react-router-dom";
import { AuthProvider } from "./auth/AuthContext";
import ProtectedRoute from "./auth/ProtectedRoute";
import Login from "./pages/Login";
import BooksList from "./pages/Books/List";
import BookCreate from "./pages/Books/Create";
import BookEdit from "./pages/Books/Edit";
import ChaptersList from "./pages/Chapters/List";
import ChapterCreate from "./pages/Chapters/Create";
import QuizBuilder from "./pages/Quizzes/QuizBuilder";

export default function App() {
  return (
    <AuthProvider>
      <BrowserRouter>
        <Routes>
          <Route path="/login" element={<Login />} />
          <Route path="/" element={<Navigate to="/books" replace />} />
          <Route path="/books" element={
            <ProtectedRoute><BooksList /></ProtectedRoute>
          } />
          <Route path="/books/new" element={
            <ProtectedRoute><BookCreate /></ProtectedRoute>
          } />
          <Route path="/books/:id" element={
            <ProtectedRoute><BookEdit /></ProtectedRoute>
          } />
          <Route path="/chapters" element={
            <ProtectedRoute><ChaptersList /></ProtectedRoute>
          } />
          <Route path="/chapters/new" element={
            <ProtectedRoute><ChapterCreate /></ProtectedRoute>
          } />
          <Route path="/quizzes/:chapterId" element={
            <ProtectedRoute><QuizBuilder /></ProtectedRoute>
          } />
        </Routes>
      </BrowserRouter>
    </AuthProvider>
  );
}