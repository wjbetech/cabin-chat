import { BrowserRouter, Navigate, Route, Routes } from "react-router-dom";
import Navbar from "./components/Navbar";
import ChatLayout from "./layouts/ChatLayout";
import PublicLayout from "./layouts/PublicLayout";
import ChatPage from "./pages/ChatPage";
import LandingPage from "./pages/LandingPage";
import LoginPage from "./pages/LoginPage";
import SignupPage from "./pages/SignupPage";

function App() {
  return (
    <BrowserRouter>
      <div className="relative isolate min-h-screen overflow-hidden text-[#5d3b20]">
        <div className="pointer-events-none absolute inset-0 bg-[radial-gradient(circle_at_top,rgba(255,255,255,0.55),transparent_38%),linear-gradient(180deg,rgba(255,255,255,0.15),transparent_24%)]" />
        <div className="pointer-events-none absolute -left-20 top-20 h-72 w-72 rounded-full bg-[#fff7cc]/60 blur-3xl" />
        <div className="pointer-events-none absolute right-0 top-16 h-96 w-96 rounded-full bg-[#c8f4ff]/50 blur-3xl" />
        <Navbar />
        <div className="relative">
          <Routes>
            <Route element={<PublicLayout />}>
              <Route path="/" element={<LandingPage />} />
              <Route path="/login" element={<LoginPage />} />
              <Route path="/signup" element={<SignupPage />} />
            </Route>
            <Route element={<ChatLayout />}>
              <Route path="/chat" element={<ChatPage />} />
            </Route>
            <Route path="*" element={<Navigate to="/" replace />} />
          </Routes>
        </div>
      </div>
    </BrowserRouter>
  );
}

export default App;
