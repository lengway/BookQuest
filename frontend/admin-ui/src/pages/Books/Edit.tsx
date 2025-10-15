import { useEffect, useState } from "react";
import { useParams } from "react-router-dom";
import { api } from "../../api/axios";
import BookForm, { BookFormValues } from "../../components/BookForm";

export default function BookEdit() {
  const { id } = useParams();
  const [data, setData] = useState<BookFormValues>();
  useEffect(()=>{ api.get(`/books/${id}`).then(r=>setData(r.data)); },[id]);

  async function onSubmit(v: BookFormValues) {
    await api.put(`/books/${id}`, v);
    alert("Saved");
  }
  if (!data) return null;
  return <BookForm defaultValues={data} onSubmit={onSubmit} />;
}