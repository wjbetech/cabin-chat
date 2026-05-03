export default function ChatPage() {
  return (
    <div className="space-y-6">
      <div className="rounded-3xl border-2 border-[#8b5e3c] bg-[#fff7e8] px-5 py-4 shadow-[0_4px_0_#8b5e3c]">
        <p className="text-sm font-extrabold uppercase tracking-[0.2em] text-[#8b5e3c]">chat canvas</p>
        <h1 className="mt-2 font-display text-3xl lowercase text-[#5d3b20]">select a room</h1>
      </div>

      <div className="space-y-4 rounded-[1.75rem] border-2 border-[#8b5e3c] bg-white/85 p-4 shadow-[0_8px_0_#8b5e3c]">
        {[
          {
            name: "maple guide",
            message: "welcome to cabin chat. the chat panel will become the live message board.",
            align: "start",
            tone: "bg-[#dff6ff]"
          },
          {
            name: "pinetree pal",
            message: "this layout is aiming for a game-like town board with soft colors and round frames.",
            align: "end",
            tone: "bg-[#fff2c7]"
          },
          {
            name: "campfire buddy",
            message: "next step: replace these sample bubbles with real room and message data.",
            align: "start",
            tone: "bg-[#ffe4ef]"
          }
        ].map((item) => (
          <div key={item.name} className={`flex ${item.align === "end" ? "justify-end" : "justify-start"}`}>
            <div
              className={`max-w-xl rounded-3xl border-2 border-[#8b5e3c] ${item.tone} px-4 py-3 shadow-[0_4px_0_#8b5e3c]`}>
              <p className="text-xs font-extrabold uppercase tracking-[0.2em] text-[#8b5e3c]">{item.name}</p>
              <p className="mt-2 text-[#5d3b20]">{item.message}</p>
            </div>
          </div>
        ))}

        <div className="mt-2 flex items-center gap-3 rounded-full border-2 border-[#8b5e3c] bg-[#fffdf8] px-4 py-3 shadow-[0_4px_0_#8b5e3c]">
          <span className="text-xl">💬</span>
          <input
            className="flex-1 bg-transparent text-[#5d3b20] outline-none placeholder:text-[#a07a55]"
            placeholder="type a message..."
          />
          <button className="rounded-full border-2 border-[#8b5e3c] bg-[#ffd86e] px-4 py-2 font-extrabold uppercase tracking-wide text-[#5d3b20] shadow-[0_3px_0_#8b5e3c]">
            send
          </button>
        </div>
      </div>
    </div>
  );
}
