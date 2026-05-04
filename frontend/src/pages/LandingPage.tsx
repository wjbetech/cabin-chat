import { Link } from "react-router-dom";
import PixelCabinScene from "../components/PixelCabinScene";

export default function LandingPage() {
  return (
    <section className="grid w-full gap-4 px-3 sm:gap-8 sm:px-8 lg:grid-cols-[minmax(0,0.95fr)_minmax(320px,1.05fr)] lg:items-center lg:px-16 xl:px-20 sm:scale-75 md:scale-90 lg:scale-100">
      <div className="space-y-3 sm:space-y-6">
        <div className="rounded-3xl border-4 border-[#8b5e3c] bg-[linear-gradient(180deg,#fffef7,#fff3d8)] p-3 shadow-[0_8px_0_#8b5e3c] sm:rounded-4xl sm:p-6 sm:shadow-[0_14px_0_#8b5e3c]">
          <h1 className="font-bold text-xl lowercase leading-none text-[#5d3b20] sm:text-3xl">cabin chat</h1>
        </div>

        <div className="flex flex-wrap gap-2 sm:gap-3">
          <Link
            to="/signup"
            className="rounded-full border-2 border-[#8b5e3c] bg-[#ffd86e] px-3 py-1.5 text-xs font-extrabold uppercase tracking-wide text-[#5d3b20] shadow-[0_3px_0_#8b5e3c] transition-transform hover:-translate-y-0.5 sm:px-6 sm:py-3 sm:text-base sm:shadow-[0_5px_0_#8b5e3c]">
            create account
          </Link>
          <Link
            to="/login"
            className="rounded-full border-2 border-[#8b5e3c] bg-[#fff8ea] px-3 py-1.5 text-xs font-extrabold uppercase tracking-wide text-[#5d3b20] shadow-[0_3px_0_#8b5e3c] transition-transform hover:-translate-y-0.5 sm:px-6 sm:py-3 sm:text-base sm:shadow-[0_5px_0_#8b5e3c]">
            log in
          </Link>
        </div>
      </div>

      <div className="mx-auto w-[82%] min-w-0 sm:w-full">
        <PixelCabinScene />
      </div>
    </section>
  );
}
