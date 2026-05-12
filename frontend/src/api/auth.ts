import { authHeaders } from "./headers";

export type SignupResponse = {
  token: string;
  username: string;
  userId: string;
};

export type LoginResponse = {
  token: string;
  username: string;
  userId: string;
};

export type SessionResponse = {
  userId: string;
  username: string;
  createdAt: string;
  status: string;
  lastSeenAt: string;
  avatarUrl: string | null;
};

export async function getSession(): Promise<SessionResponse> {
  const res = await fetch("/api/session", {
    method: "GET",
    headers: authHeaders()
  });

  if (!res.ok) {
    throw new Error(`session check failed: ${res.status}`);
  }

  return res.json();
}

export async function signup(username: string, password: string): Promise<SignupResponse> {
  const res = await fetch("/api/signup", {
    method: "POST",
    headers: {
      "Content-Type": "application/json"
    },
    body: JSON.stringify({ username, password })
  });

  if (!res.ok) {
    const body = await res.json().catch(() => ({}));
    throw new Error(body.error ?? `signup failed: ${res.status}`);
  }

  return res.json();
}

export async function login(username: string, password: string): Promise<LoginResponse> {
  const res = await fetch("/api/login", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ username, password })
  });

  if (!res.ok) {
    const body = await res.json().catch(() => ({}));
    throw new Error(body.error ?? `login failed: ${res.status}`);
  }

  return res.json();
}
