import { Navigate } from "react-router-dom";

import { isAuthenticated } from "../auth";

export function AuthGuard({ children }: { children: React.ReactNode }) {
  if (!isAuthenticated()) {
    return <Navigate to="/login" replace />;
  }
  return <>{children}</>;
}
