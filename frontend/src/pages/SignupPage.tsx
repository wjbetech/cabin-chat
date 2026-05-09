import { Link } from "react-router-dom";
import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { signup } from "../api/auth";
import { useAuth } from "../context/useAuth";

export default function SignupPage() {
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);
  const navigate = useNavigate();
  const { setAuth } = useAuth();

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    setError(null);
    setLoading(true);
    try {
      const { token, username: signedInUsername } = await signup(username, password);
      setAuth(token, signedInUsername);
      navigate("/chat");
    } catch (error) {
      setError(error instanceof Error ? error.message : "An unknown error occurred!");
    } finally {
      setLoading(false);
    }
  }

  return (
    <form
      className="w-full max-w-md rounded-4xl border-4 border-[#8b5e3c] bg-[linear-gradient(180deg,#fffdf6,#fff1d0)] p-4 shadow-[0_14px_0_#8b5e3c] mx-auto"
      onSubmit={handleSubmit}>
      <div className="rounded-3xl border-2 border-[#8b5e3c] bg-white/85 p-6">
        <div>
          <h1 className="mb-2 font-display text-4xl  text-[#5d3b20]">Sign Up</h1>
        </div>

        <div className="flex flex-col space-y-6 mt-4 mb-4">
          <label className="form-control">
            <span className="font-semibold text-[#5d3b20]">username</span>
            <input
              className="input input-bordered w-full rounded-2xl border-2 border-[#8b5e3c] bg-[#fffdf8]"
              placeholder="user"
              value={username}
              onChange={(e) => setUsername(e.target.value)}
            />
          </label>

          <label className="form-control">
            <span className="mb-4 font-semibold text-[#5d3b20]">password</span>
            <input
              className="input input-bordered w-full rounded-2xl border-2 border-[#8b5e3c] bg-[#fffdf8]"
              type="password"
              placeholder="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
            />
          </label>
        </div>

        {error && (
          <div className="flex justify-between w-full rounded-lg bg-red-100 p-4">
            <p className="text-sm text-red-700"></p>
            {error}{" "}
            <span className="cursor-pointer" onClick={() => setError(null)}>
              x
            </span>
          </div>
        )}

        <div className="mt-6 flex flex-col items-stretch gap-3">
          <button
            type="submit"
            disabled={loading}
            className="rounded-full border-2 border-[#8b5e3c] bg-[#ffd86e] px-6 py-3 font-extrabold uppercase tracking-wide text-[#5d3b20] shadow-[0_5px_0_#8b5e3c] transition-transform hover:-translate-y-0.5 cursor-pointer">
            {loading ? "Creating account..." : "Sign up"}
          </button>
          <p className="text-sm text-[#6f4a29] mt-2">
            Already have an account?{" "}
            <Link to="/login" className="font-bold text-[#8b5e3c] underline decoration-2 underline-offset-4">
              log in
            </Link>
          </p>
        </div>
      </div>
    </form>
  );
}
