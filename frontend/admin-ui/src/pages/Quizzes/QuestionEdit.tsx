import React, { useEffect, useState } from "react";
import { useParams, Link, useNavigate } from "react-router-dom";
import { api } from "../../api/axios";

type Option = { id:number; text:string; is_correct?:boolean; order_index?:number; match_key?:string };
type Question = { id:number; type:"single_choice"|"multi_choice"|"ordering"|"matching"; text:string; order_index:number; score:number; options:Option[] };

export default function QuestionEdit() {
  const { questionId, chapterId } = useParams();
  const nav = useNavigate();
  const [q, setQ] = useState<Question | null>(null);
  const [loading, setLoading] = useState(true);

  async function load() {
    // У нас нет get by id эндпоинта для вопроса, поэтому берём квиз по главе и вытаскиваем вопрос
    const { data: quiz } = await api.get(`/quizzes/by-chapter/${chapterId}`);
    const found = (quiz.questions as Question[]).find((x)=> String(x.id) === String(questionId)) || null;
    setQ(found);
    setLoading(false);
  }

  useEffect(()=>{ load(); /*eslint-disable-next-line*/ },[questionId, chapterId]);

  async function saveQuestion() {
    if (!q) return;
    await api.put(`/quizzes/questions/${q.id}`, {
      type: q.type,
      text: q.text,
      order_index: q.order_index,
      score: q.score,
    });
    alert("Question saved");
  }

  async function saveOption(o: Option) {
    await api.put(`/quizzes/options/${o.id}`, {
      text: o.text,
      is_correct: !!o.is_correct,
      order_index: o.order_index,
      match_key: o.match_key,
    });
  }

  async function deleteOption(o: Option) {
    await api.delete(`/quizzes/options/${o.id}`);
    setQ((prev)=> prev ? { ...prev, options: prev.options.filter(x=>x.id!==o.id)} : prev);
  }

  async function addOption() {
    if (!q) return;
    const { data } = await api.post(`/quizzes/questions/${q.id}/options`, {
      text: "New option",
      is_correct: false,
    });
    setQ({ ...q, options: [...q.options, data]});
  }

  if (loading) return <div>Loading...</div>;
  if (!q) return <div>Question not found. <Link to={`/quizzes/${chapterId}`}>Back</Link></div>;

  return (
    <div>
      <Link to={`/quizzes/${chapterId}`}>Back to Quiz</Link>
      <h2>Edit Question #{q.order_index+1}</h2>

      <div style={{display:"grid", gap:8, maxWidth:480}}>
        <label>
          Type:
          <select value={q.type} onChange={(e)=>setQ({...q, type: e.target.value as Question["type"]})}>
            <option value="single_choice">single_choice</option>
            <option value="multi_choice">multi_choice</option>
            <option value="ordering">ordering</option>
            <option value="matching">matching</option>
          </select>
        </label>

        <label>
          Text:
          <input value={q.text} onChange={(e)=>setQ({...q, text: e.target.value})}/>
        </label>

        <label>
          Order:
          <input type="number" value={q.order_index} onChange={(e)=>setQ({...q, order_index: Number(e.target.value)})}/>
        </label>

        <label>
          Score:
          <input type="number" value={q.score} onChange={(e)=>setQ({...q, score: Number(e.target.value)})}/>
        </label>

        <button onClick={saveQuestion}>Save Question</button>
      </div>

      <h3>Options</h3>
      <button onClick={addOption}>+ Add option</button>
      <ul>
        {q.options.map((o)=>(
          <li key={o.id} style={{margin:"8px 0"}}>
            <input
              value={o.text}
              onChange={(e)=>setQ(prev=> prev ? {...prev, options: prev.options.map(x=> x.id===o.id ? {...x, text: e.target.value} : x)} : prev)}
              style={{ width: 260 }}
            />
            {" "}
            <label>
              Correct
              <input
                type="checkbox"
                checked={!!o.is_correct}
                onChange={(e)=>setQ(prev=> prev ? {...prev, options: prev.options.map(x=> x.id===o.id ? {...x, is_correct: e.target.checked} : x)} : prev)}
              />
            </label>
            {" "}
            <label>
              Order
              <input
                type="number"
                value={o.order_index ?? 0}
                onChange={(e)=>setQ(prev=> prev ? {...prev, options: prev.options.map(x=> x.id===o.id ? {...x, order_index: Number(e.target.value)} : x)} : prev)}
                style={{ width: 80 }}
              />
            </label>
            {" "}
            <label>
              Match key
              <input
                value={o.match_key ?? ""}
                onChange={(e)=>setQ(prev=> prev ? {...prev, options: prev.options.map(x=> x.id===o.id ? {...x, match_key: e.target.value} : x)} : prev)}
                style={{ width: 120 }}
              />
            </label>
            {" "}
            <button onClick={()=>saveOption(q.options.find(x=>x.id===o.id)!)}>Save</button>
            {" "}
            <button onClick={()=>deleteOption(o)}>Delete</button>
          </li>
        ))}
      </ul>
    </div>
  );
}