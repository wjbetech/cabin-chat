export type SignupResponse = {
  token: string;
  user: {
    id: string;
    username: string;
  };
};

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
