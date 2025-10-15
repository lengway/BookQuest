import React, { useEffect, useState } from "react";
import { Link } from "react-router-dom";
import { api } from "../../api/axios";

type Chapter = { id: number; book_id: number; chapter_number: number; title: string };

export default function QuizzesList() {
  const [chapters, setChapters] = useState<Chapter[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    api.get("/chapters").then(r => setChapters(r.data)).finally(() => setLoading(false));
  }, []);

  if (loading) return <div>Loading...</div>;
  return (
    <div>
      <Link to="/">Home</Link>
      <h2>Quizzes by Chapter</h2>
      <ul>
        {chapters.map(c => (
          <li key={c.id}>
            <span>Book {c.book_id} · #{c.chapter_number} — {c.title}</span>{" "}
            <Link to={`/quizzes/${c.id}`}>Open quiz</Link>
          </li>
        ))}
      </ul>
    </div>
  );
}