import { Link } from "react-router-dom";

export default function LandingPage() {
  return (
    <section className="grid gap-8 lg:grid-cols-[minmax(0,1.2fr)_minmax(320px,0.8fr)] lg:items-center">
      <div className="rounded-4xl border-4 border-[#8b5e3c] bg-[linear-gradient(180deg,#fffef7,#fff3d8)] p-4 shadow-[0_14px_0_#8b5e3c]">
        <h1 className="text-2xl">cabin chat</h1>
      </div>
      <div className="space-y-6">
        <div className="flex flex-wrap gap-3">
          <Link
            to="/signup"
            className="rounded-full border-2 border-[#8b5e3c] bg-[#ffd86e] px-6 py-3 font-extrabold uppercase tracking-wide text-[#5d3b20] shadow-[0_5px_0_#8b5e3c] transition-transform hover:-translate-y-0.5">
            create account
          </Link>
          <Link
            to="/login"
            className="rounded-full border-2 border-[#8b5e3c] bg-[#fff8ea] px-6 py-3 font-extrabold uppercase tracking-wide text-[#5d3b20] shadow-[0_5px_0_#8b5e3c] transition-transform hover:-translate-y-0.5">
            log in
          </Link>
        </div>
      </div>
    </section>
  );
}
