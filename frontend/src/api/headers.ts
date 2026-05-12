export function authHeaders(): HeadersInit {
  const token = localStorage.getItem("cabin_token");
  if (!token) {
    return {
      "Content-Type": "application/json"
    };
  } else {
    return {
      "Content-Type": "application/json",
      Authorization: `Bearer ${token}`
    };
  }
}
