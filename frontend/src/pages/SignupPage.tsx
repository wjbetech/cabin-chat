import { Link } from "react-router-dom";

export default function SignupPage() {
  return (
    <div className="w-full max-w-md rounded-4xl border-4 border-[#8b5e3c] bg-[linear-gradient(180deg,#fffdf6,#fff1d0)] p-4 shadow-[0_14px_0_#8b5e3c] mx-auto">
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
            />
          </label>

          <label className="form-control">
            <span className="mb-4 font-semibold text-[#5d3b20]">password</span>
            <input
              className="input input-bordered w-full rounded-2xl border-2 border-[#8b5e3c] bg-[#fffdf8]"
              type="password"
              placeholder="password"
            />
          </label>
        </div>

        <div className="mt-6 flex flex-col items-stretch gap-3">
          <button className="rounded-full border-2 border-[#8b5e3c] bg-[#ffd86e] px-6 py-3 font-extrabold uppercase tracking-wide text-[#5d3b20] shadow-[0_5px_0_#8b5e3c] transition-transform hover:-translate-y-0.5 cursor-pointer">
            create account
          </button>
          <p className="text-sm text-[#6f4a29] mt-2">
            Already have an account?{" "}
            <Link to="/login" className="font-bold text-[#8b5e3c] underline decoration-2 underline-offset-4">
              log in
            </Link>
          </p>
        </div>
      </div>
    </div>
  );
}
