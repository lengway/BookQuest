import { useForm } from "react-hook-form";
import { api } from "../../api/axios";

export default function ChapterCreate() {
  const { register, handleSubmit } = useForm<{book_id:number; chapter_number:number; title:string; content:string}>();
  async function onSubmit(v:any) {
    v.book_id = Number(v.book_id);
    v.chapter_number = Number(v.chapter_number);
    await api.post("/chapters", v);
    alert("Created");
  }
  return (
    <form onSubmit={handleSubmit(onSubmit)}>
      <input placeholder="Book ID" type="number" {...register("book_id", {valueAsNumber:true})}/>
      <input placeholder="Chapter #" type="number" {...register("chapter_number", {valueAsNumber:true})}/>
      <input placeholder="Title" {...register("title", {required:true})}/>
      <textarea placeholder="Content" {...register("content", {required:true})}/>
      <button type="submit">Create</button>
    </form>
  );
}