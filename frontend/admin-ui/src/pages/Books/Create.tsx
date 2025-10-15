import { api } from "../../api/axios";
import { useNavigate } from "react-router-dom";
import BookForm, { BookFormValues } from "../../components/BookForm";

export default function BookCreate() {
  const nav = useNavigate();
  async function onSubmit(v: BookFormValues) {
    const { data } = await api.post("/books", v);
    nav(`/books/${data.id}`);
  }
  return <BookForm onSubmit={onSubmit} />;
}