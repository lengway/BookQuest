import React, { createContext, useContext, useState, type ReactNode } from "react";
import { api } from "../api/axios";

type User = { id: number; username: string; is_superuser: boolean };
type AuthCtx = {
  user: User | null;
  login: (email: string, password: string) => Promise<void>;
  logout: () => void;
};
const Ctx = createContext<AuthCtx>(null as any);

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [user, setUser] = useState<User | null>(null);

  async function login(email: string, password: string) {
    const { data } = await api.post("/auth/login", { email, password });
    localStorage.setItem("access_token", data.access_token);
    const me = await api.get("/auth/me");
    setUser(me.data);
  }
  function logout() {
    localStorage.removeItem("access_token");
    setUser(null);
  }
  return <Ctx.Provider value={{ user, login, logout }}>{children}</Ctx.Provider>;
}
export const useAuth = () => useContext(Ctx);