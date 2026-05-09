import { createContext, useContext } from "react";

export type AuthUser = { username: string };

export type AuthContextType = {
  user: AuthUser | null;
  setAuth: (token: string, username: string) => void;
  logout: () => void;
};

export const AuthContext = createContext<AuthContextType | null>(null);

export function useAuth() {
  const ctx = useContext(AuthContext);
  if (!ctx) {
    throw new Error("useAuth must be used within an AuthProvider");
  }
  return ctx;
}
