import React, { useEffect, useState } from "react";
import { Link } from "react-router-dom";
import { api } from "../../api/axios";

type Chapter = {
    id: number;
    book_id: number;
    chapter_number: number;
    title: string;
};

export default function ChaptersList() {
    const [items, setItems] = useState<Chapter[]>([]);
    const [bookId, setBookId] = useState<string>("");

    async function load() {
        const params: Record<string, string | number> = {};
        if (bookId) params.book_id = Number(bookId);
        const { data } = await api.get("/chapters", { params });
        setItems(data);
    }

    useEffect(() => {
        load();
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, []);

    return (
        <div>
            <Link to="/">Home</Link>
            <h2>Chapters</h2>
            <div style={{ display: "flex", gap: 8, alignItems: "center" }}>
                <input
                    placeholder="Filter by Book ID"
                    value={bookId}
                    onChange={(e) => setBookId(e.target.value)}
                    style={{ width: 160 }}
                />
                <button onClick={load}>Apply</button>
                <Link to="/chapters/new">Create</Link>
            </div>
            <ul>
                {items.map((c) => (
                    <li key={c.id}>
                        <span>Book {c.book_id} · #{c.chapter_number} — {c.title}</span>
                    </li>
                ))}
            </ul>
        </div>
    );
}


