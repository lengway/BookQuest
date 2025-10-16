import { useForm } from "react-hook-form";

export type ChapterFormValues = {
    book_id: number; 
    chapter_number: number; 
    title: string; 
    content: string; 
    estimated_reading_time: number;
};

export default function ChapterForm({ defaultValues, onSubmit }:{
    defaultValues?: Partial<ChapterFormValues>;
    onSubmit: (v: ChapterFormValues)=>void;
}) {
    const { register, handleSubmit } = useForm<ChapterFormValues>({ defaultValues });
    return (
        <form onSubmit={handleSubmit(onSubmit)}>
            <label>Book ID </label><br/><input placeholder="Book ID" type="number" {...register("book_id", {valueAsNumber:true})}/><br/>
            <label>Chapter # </label><br/><input placeholder="Chapter #" type="number" {...register("chapter_number", {valueAsNumber:true})}/><br/>
            <label>Title </label><br/><input placeholder="Title" {...register("title", {required:true})}/><br/>
            <label>Content </label><br/><textarea placeholder="Content" {...register("content", {required:true})}/><br/>
            <label>Estimated Reading Time </label><br/><input placeholder="Estimated Reading Time" type="number" {...register("estimated_reading_time", {valueAsNumber:true})}/><br/>
            <button type="submit">Save</button>
        </form>
    );
}