import { createContext, useContext, useState } from "react";

type AuthUser = { username: string };

type AuthContextType = {
  user: AuthUser | null;
  setAuth: (token: string, username: string) => void;
  logout: () => void;
};

const AuthContext = createContext<AuthContextType | null>(null);

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

export function useAuth() {
  const ctx = useContext(AuthContext);
  if (!ctx) {
    throw new Error("useAuth must be used within an AuthProvider");
  }
  return ctx;
}
