import { useState, useEffect } from "react";
import { AuthContext, type AuthUser } from "./useAuth";
import { getSession } from "../api/auth";

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [user, setUser] = useState<AuthUser | null>(() => {
    const username = localStorage.getItem("cabin_username");
    return username ? { username } : null;
  });
  const [ready, setReady] = useState(() => !localStorage.getItem("cabin_token"));

  useEffect(() => {
    const token = localStorage.getItem("cabin_token");

    if (!token) {
      return;
    }

    getSession()
      .then((session) => {
        setUser({ username: session.username });
      })
      .catch(() => {
        localStorage.removeItem("cabin_token");
        localStorage.removeItem("cabin_username");
        setUser(null);
      })
      .finally(() => {
        setReady(true);
      });
  }, []);

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

  if (!ready) {
    return null; // or a loading spinner
  }

  return <AuthContext.Provider value={{ user, setAuth, logout }}>{children}</AuthContext.Provider>;
}
