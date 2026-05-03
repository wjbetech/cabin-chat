import { Outlet } from "react-router-dom";

export default function PublicLayout() {
  return (
    <main className="mx-auto flex min-h-[calc(100vh-4.5rem)] w-full max-w-6xl items-center px-4 py-10">
      <Outlet />
    </main>
  );
}
