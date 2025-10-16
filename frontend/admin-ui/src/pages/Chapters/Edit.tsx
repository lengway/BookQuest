import { useEffect, useState } from "react";
import { useParams } from "react-router-dom";
import { api } from "../../api/axios";
import ChapterForm, { ChapterFormValues } from "../../components/ChapterForm";

export default function ChapterEdit() {
  const { id } = useParams();
  const [data, setData] = useState<ChapterFormValues>();
  useEffect(()=>{ api.get(`/chapters/${id}`).then(r=>setData(r.data)); },[id]);
  async function onSubmit(v: ChapterFormValues) {
    await api.put(`/chapters/${id}`, v);
    alert("Saved");
  }
  if (!data) return null;
  return <ChapterForm defaultValues={data} onSubmit={onSubmit} />;
}