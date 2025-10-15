import { useEffect, useState } from "react";
import { api } from "../../api/axios";
import { Link } from "react-router-dom";

type Book = { id: number; title: string; author: string };

export default function BooksList() {
  const [items, setItems] = useState<Book[]>([]);
  useEffect(() => { api.get("/books").then(r => setItems(r.data)); }, []);
  return (
    <div>
      <Link to="/">Home</Link>
      <h2>Books</h2>
      <Link to="/books/new">Create</Link>
      <ul>{items.map(b => <li key={b.id}>
        <Link to={`/books/${b.id}`}>{b.title} — {b.author}</Link>
      </li>)}</ul>
    </div>
  );
}