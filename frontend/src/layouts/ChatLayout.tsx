import { Outlet } from "react-router-dom";

export default function ChatLayout() {
  return (
    <main className="mx-auto grid min-h-[calc(100vh-4rem)] w-full max-w-7xl gap-5 px-4 py-6 lg:grid-cols-[300px_minmax(0,1fr)]">
      <aside className="rounded-4xl border-4 border-[#8b5e3c] bg-[linear-gradient(180deg,#fff8e9,#fff3cf)] p-4 shadow-[0_12px_0_#8b5e3c]">
        <div className="space-y-4">
          <div className="rounded-3xl border-2 border-[#8b5e3c] bg-white/70 p-4 shadow-[0_4px_0_#8b5e3c]">
            <p className="text-sm font-extrabold uppercase tracking-[0.2em] text-[#8b5e3c]">party board</p>
            <p className="mt-1 font-display text-2xl lowercase text-[#5d3b20]">town list</p>
          </div>

          <div className="space-y-3">
            {["henesys plaza", "sleepy camp", "fox tree hideout", "victoria outpost"].map((room) => (
              <div
                key={room}
                className="flex items-center justify-between rounded-[1.25rem] border-2 border-[#8b5e3c] bg-[#fffdf6] px-4 py-3 shadow-[0_3px_0_#8b5e3c]">
                <span className="font-semibold capitalize text-[#5d3b20]">{room}</span>
                <span className="badge border-0 bg-[#ffe08a] text-[#5d3b20]">12</span>
              </div>
            ))}
          </div>
        </div>
      </aside>

      <section className="rounded-4xl border-4 border-[#8b5e3c] bg-[linear-gradient(180deg,#fffef9,#fff4db)] p-4 shadow-[0_12px_0_#8b5e3c]">
        <Outlet />
      </section>
    </main>
  );
}
