import { NavLink, useNavigate } from "react-router-dom";
import { useAuth } from "../context/useAuth";

const linkClassName = ({ isActive }: { isActive: boolean }) =>
  [
    "btn btn-sm rounded-full border-2 border-[#8b5e3c] bg-[#fff5d7] px-4 text-sm font-extrabold uppercase tracking-wide text-[#5d3b20] shadow-[0_4px_0_#8b5e3c] transition-transform duration-150 hover:-translate-y-0.5 hover:bg-[#ffe7a8]",
    isActive ? "bg-[#ffd37a]" : ""
  ]
    .filter(Boolean)
    .join(" ");

export default function Navbar() {
  const { user, logout } = useAuth();
  const navigate = useNavigate();

  function handleLogout() {
    logout();
    navigate("/");
  }

  return (
    <div className="navbar relative z-10 border-b-4 border-[#8b5e3c] bg-[linear-gradient(180deg,#fff7e5,#ffecc7)] px-4 py-3 shadow-[0_8px_0_#8b5e3c] backdrop-blur">
      <div className="navbar-start">
        <NavLink to="/" className="flex items-center gap-3 text-2xl font-black lowercase tracking-tight text-[#5d3b20]">
          <span className="grid h-11 w-11 place-items-center rounded-2xl border-2 border-[#8b5e3c] bg-[linear-gradient(180deg,#fff5d7,#ffd77a)] text-xl shadow-[0_4px_0_#8b5e3c]">
            🌳
          </span>
          <span className="font-display text-2xl lowercase">cabin chat</span>
        </NavLink>
      </div>

      <div className="navbar-end gap-2">
        {user ? (
          <>
            <span className="mr-2">Hi, {user.username}!</span>
            <NavLink to="/" className={linkClassName}>
              home
            </NavLink>
            <NavLink to="/chat" className={linkClassName}>
              chat
            </NavLink>
            <NavLink to="/" onClick={handleLogout} className={linkClassName}>
              log out
            </NavLink>
          </>
        ) : (
          <>
            <NavLink to="/signup" className={linkClassName}>
              sign up
            </NavLink>
            <NavLink to="/login" className={linkClassName}>
              login
            </NavLink>
          </>
        )}
      </div>
    </div>
  );
}
