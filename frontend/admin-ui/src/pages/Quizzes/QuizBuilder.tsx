import { useEffect, useState } from "react";
import { useParams } from "react-router-dom";
import { api } from "../../api/axios";

type Option = { id: number; text: string; is_correct?: boolean; order_index?: number; match_key?: string };
type Question = { id: number; type: string; text: string; order_index: number; options: Option[] };
type Quiz = { id: number; chapter_id: number; title: string; description?: string; questions: Question[] };

export default function QuizBuilder() {
  const { chapterId } = useParams();
  const [quiz, setQuiz] = useState<Quiz | null>(null);

  useEffect(() => {
    api.get(`/quizzes/by-chapter/${chapterId}`).then(r => setQuiz(r.data)).catch(async () => {
      // если квиза нет — создать (нужен суперюзер)
      const created = await api.post("/quizzes", { chapter_id: Number(chapterId), title: `Quiz for chapter ${chapterId}` });
      setQuiz(created.data);
    });
  }, [chapterId]);

  async function addQuestion(type: string) {
    if (!quiz) return;
    const { data } = await api.post(`/quizzes/${quiz.id}/questions`, { type, text: "New question", order_index: quiz.questions.length });
    setQuiz({ ...quiz, questions: [...quiz.questions, data] });
  }

  async function addOption(q: Question) {
    if (!quiz) return;
    const { data } = await api.post(`/quizzes/questions/${q.id}/options`, { text: "Option", is_correct: false });
    setQuiz({
      ...quiz,
      questions: quiz.questions.map(x => x.id === q.id ? { ...x, options: [...x.options, data] } : x)
    });
  }

  if (!quiz) return <div>Loading...</div>;
  return (
    <div>
      <h2>Quiz for chapter {quiz.chapter_id}</h2>
      <button onClick={() => addQuestion("single_choice")}>+ Single</button>
      <button onClick={() => addQuestion("multi_choice")}>+ Multi</button>
      <button onClick={() => addQuestion("ordering")}>+ Ordering</button>
      <button onClick={() => addQuestion("matching")}>+ Matching</button>

      {quiz.questions.map(q => (
        // внутри map по quiz.questions:
        <div key={q.id} style={{ border: "1px solid #ddd", margin: "8px", padding: "8px" }}>
          <div>Q#{q.order_index + 1} [{q.type}] {q.text}</div>
          <button onClick={() => addOption(q)}>+ Option</button>
          <a href={`/quizzes/${quiz.chapter_id}/questions/${q.id}`} style={{ marginLeft: 8 }}>Edit</a>
          <ul>{q.options.map(o => <li key={o.id}>{o.text} {o.is_correct ? "(correct)" : ""}</li>)}</ul>
        </div>
      ))}
    </div>
  );
}