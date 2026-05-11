import { useState } from "react";
import { Link, useNavigate } from "react-router-dom";
import { useAuth } from "../context/useAuth";
import { login } from "../api/auth";

export default function LoginPage() {
  const { setAuth } = useAuth();
  const navigate = useNavigate();
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    setError(null);
    setLoading(true);
    try {
      const { token, username: loggedInUsername } = await login(username, password);
      setAuth(token, loggedInUsername);
      navigate("/chat");
    } catch (error) {
      setError(error instanceof Error ? error.message : "An error occurred while logging in!");
    } finally {
      setLoading(false);
    }
  }

  return (
    <div className="w-full max-w-md rounded-4xl border-4 border-[#8b5e3c] bg-[linear-gradient(180deg,#fffdf6,#fff1d0)] p-4 shadow-[0_14px_0_#8b5e3c] mx-auto">
      <div className="rounded-3xl border-2 border-[#8b5e3c] bg-white/85 p-6">
        <div>
          <h1 className="mb-2 font-display text-4xl text-[#5d3b20]">Log In</h1>
        </div>

        <form onSubmit={handleSubmit}>
          <div className="flex flex-col space-y-6 mt-4 mb-4">
            <label className="form-control">
              <span className="font-semibold text-[#5d3b20]">username</span>
              <input
                className="input input-bordered w-full rounded-2xl border-2 border-[#8b5e3c] bg-[#fffdf8]"
                placeholder="will"
                value={username}
                onChange={(e) => setUsername(e.target.value)}
              />
            </label>

            <label className="form-control">
              <span className="mb-4 font-semibold text-[#5d3b20]">password</span>
              <input
                className="input input-bordered w-full rounded-2xl border-2 border-[#8b5e3c] bg-[#fffdf8]"
                type="password"
                placeholder="••••••••"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
              />
            </label>
          </div>
          {error && (
            <div className="rounded-2xl border-2 border-[#e55353] bg-[#ffe5e5] px-4 py-3 text-sm text-[#a23b3b] mb-4">
              {error}
            </div>
          )}

          <div className="mt-6 flex flex-col items-stretch gap-3">
            <button
              type="submit"
              className="rounded-full border-2 border-[#8b5e3c] bg-[#ffd86e] px-6 py-3 font-extrabold uppercase tracking-wide text-[#5d3b20] shadow-[0_5px_0_#8b5e3c] transition-transform hover:-translate-y-0.5 cursor-pointer"
              disabled={loading}>
              log in
            </button>
            <p className="text-sm text-[#6f4a29] mt-2">
              Don't have an account?{" "}
              <Link to="/signup" className="font-bold text-[#8b5e3c] underline decoration-2 underline-offset-4">
                sign up
              </Link>
            </p>
          </div>
        </form>
      </div>
    </div>
  );
}
