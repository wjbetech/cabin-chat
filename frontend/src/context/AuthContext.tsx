import { useState } from "react";
import { AuthContext, type AuthUser } from "./useAuth";

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [user, setUser] = useState<AuthUser | null>(() => {
    const username = localStorage.getItem("cabin_username");
    return username ? { username } : null;
  });

  function setAuth(token: string, username: string) {
    localStorage.setItem("cabin_token", token);
    localStorage.setItem("cabin_username", username);
    setUser({ username });
  }

  function logout() {
    localStorage.removeItem("cabin_token");
    localStorage.removeItem("cabin_username");
    setUser(null);
  }

  return <AuthContext.Provider value={{ user, setAuth, logout }}>{children}</AuthContext.Provider>;
}
