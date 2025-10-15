import { useForm } from "react-hook-form";

export type BookFormValues = {
  title: string; author: string; description?: string;
  cover_image_url?: string; genre?: string; difficulty?: string;
  total_chapters?: number; estimated_reading_time?: number;
  language?: string; chapter_xp?: number; completion_xp?: number;
};

export default function BookForm({ defaultValues, onSubmit }:{
  defaultValues?: Partial<BookFormValues>;
  onSubmit: (v: BookFormValues)=>void;
}) {
  const { register, handleSubmit } = useForm<BookFormValues>({ defaultValues });
  return (
    <form onSubmit={handleSubmit(onSubmit)}>
      <input placeholder="Title" {...register("title", { required: true })} />
      <input placeholder="Author" {...register("author", { required: true })} />
      <input placeholder="Genre" {...register("genre")} />
      <input placeholder="Language" defaultValue="ru" {...register("language")} />
      <input placeholder="Chapter XP" type="number" {...register("chapter_xp", {valueAsNumber:true})} />
      <textarea placeholder="Description" {...register("description")} />
      <button type="submit">Save</button>
    </form>
  );
}