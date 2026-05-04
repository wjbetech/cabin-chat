import { Link } from "react-router-dom";
import PixelCabinScene from "../components/PixelCabinScene";
import cat1 from "../assets/cat-1.jpg";
import cat2 from "../assets/cat-2.jpg";

export default function LandingPage() {
  return (
    <section className="grid w-full gap-4 px-3 sm:gap-8 sm:px-8 lg:grid-cols-[minmax(0,0.95fr)_minmax(320px,1.05fr)] lg:items-center lg:px-16 xl:px-20 sm:scale-80 md:scale-100 justify-center">
      <div className="space-y-3 sm:space-y-6 max-w-100">
        <div className="chat chat-start">
          <div className="chat-image avatar">
            <div className="w-10 rounded-full">
              <img src={cat1} alt="Cat avatar" className="object-cover [image-rendering:pixelated]" />
            </div>
          </div>
          <div className="chat-bubble chat-bubble-primary">hey, want to join our cabin?</div>
        </div>
        <div className="chat chat-end">
          <div className="chat-image avatar">
            <div className="w-10 rounded-full">
              <img src={cat2} alt="Cat avatar flipped" className="object-cover [image-rendering:pixelated]" />
            </div>
          </div>
          <div className="chat-bubble chat-bubble-secondary">already here! 🐾</div>
        </div>

        <div className="flex flex-wrap gap-2 sm:gap-3 justify-between">
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
