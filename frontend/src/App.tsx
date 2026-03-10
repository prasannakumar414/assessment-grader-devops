import { Navigate, Route, Routes, useOutletContext } from "react-router-dom";

import { AuthGuard } from "./components/AuthGuard";
import { Layout } from "./components/Layout";
import { AddStudentPage } from "./pages/AddStudent";
import { LoginPage } from "./pages/Login";
import { RegistrationRequestsPage } from "./pages/RegistrationRequests";
import { StudentListPage } from "./pages/StudentList";
import { StudentProfilePage } from "./pages/StudentProfile";

interface OutletCtx {
  registrationVersion: number;
  stageVersion: number;
}

function StudentListWrapper() {
  const { stageVersion } = useOutletContext<OutletCtx>();
  return <StudentListPage stageVersion={stageVersion} />;
}

function RegistrationRequestsWrapper() {
  const { registrationVersion } = useOutletContext<OutletCtx>();
  return <RegistrationRequestsPage registrationVersion={registrationVersion} />;
}

function App() {
  return (
    <Routes>
      <Route path="/login" element={<LoginPage />} />
      <Route
        element={
          <AuthGuard>
            <Layout />
          </AuthGuard>
        }
        path="/"
      >
        <Route element={<StudentListWrapper />} index />
        <Route element={<RegistrationRequestsWrapper />} path="registrations" />
        <Route element={<AddStudentPage />} path="add" />
        <Route element={<StudentProfilePage />} path="students/:id" />
        <Route element={<Navigate to="/" replace />} path="*" />
      </Route>
    </Routes>
  );
}

export default App;
