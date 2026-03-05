import { NavLink, Outlet, useNavigate } from "react-router-dom";

import { clearToken } from "../auth";
import { useSSE } from "../hooks/useSSE";
import { AllCompleteModal } from "./AllCompleteModal";
import { StageCompleteModal } from "./StageCompleteModal";

const links = [
  { to: "/", label: "Students" },
  { to: "/registrations", label: "Registrations" },
  { to: "/add", label: "Add Student" },
];

export function Layout() {
  const navigate = useNavigate();
  const { connected, currentEvent, dismissCurrent, registrationVersion } = useSSE();

  function handleLogout() {
    clearToken();
    navigate("/login", { replace: true });
  }

  return (
    <div className="min-h-screen bg-slate-50 text-slate-900">
      <header className="border-b border-slate-200 bg-white">
        <div className="mx-auto flex max-w-6xl items-center justify-between px-4 py-3">
          <div className="flex items-center gap-3">
            <h1 className="text-lg font-bold">DevOps Assessment Grader</h1>
            {!connected && (
              <span className="rounded bg-yellow-100 px-2 py-0.5 text-xs font-medium text-yellow-800">
                Reconnecting...
              </span>
            )}
          </div>
          <div className="flex items-center gap-3">
            <nav className="flex gap-3">
              {links.map((link) => (
                <NavLink
                  key={link.to}
                  to={link.to}
                  className={({ isActive }) =>
                    `rounded px-3 py-2 text-sm font-medium ${
                      isActive ? "bg-slate-900 text-white" : "text-slate-700 hover:bg-slate-100"
                    }`
                  }
                >
                  {link.label}
                </NavLink>
              ))}
            </nav>
            <button
              onClick={handleLogout}
              className="rounded px-3 py-2 text-sm font-medium text-red-600 hover:bg-red-50"
            >
              Logout
            </button>
          </div>
        </div>
      </header>

      <main className="mx-auto max-w-6xl px-4 py-6">
        <Outlet context={{ registrationVersion }} />
      </main>

      {currentEvent?.type === "stage_complete" && (
        <StageCompleteModal event={currentEvent.data} onDismiss={dismissCurrent} />
      )}
      {currentEvent?.type === "all_complete" && (
        <AllCompleteModal event={currentEvent.data} onDismiss={dismissCurrent} />
      )}
    </div>
  );
}
