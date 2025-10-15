import { useForm } from "react-hook-form";
import { useAuth } from "../auth/AuthContext";

export default function Login() {
  const { register, handleSubmit } = useForm<{email:string; password:string}>();
  const { login } = useAuth();
  return (
    <form onSubmit={handleSubmit(({email,password})=>login(email,password))}>
      <input placeholder="Email" {...register("email")} />
      <input placeholder="Password" type="password" {...register("password")} />
      <button type="submit">Login</button>
    </form>
  );
}