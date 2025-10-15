import React from "react";
import { Link } from "react-router-dom";

export default function Home() {
  return (
    <div style={{ display: "flex", gap: 16, flexWrap: "wrap" }}>
      <Link to="/books">Books</Link>
      <Link to="/chapters">Chapters</Link>
      <Link to="/quizzes/">Quizzes (by chapter)</Link>
    </div>
  );
}