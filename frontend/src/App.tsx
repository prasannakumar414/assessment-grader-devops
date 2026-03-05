import { Navigate, Route, Routes, useOutletContext } from "react-router-dom";

import { Layout } from "./components/Layout";
import { AddStudentPage } from "./pages/AddStudent";
import { RegistrationRequestsPage } from "./pages/RegistrationRequests";
import { StudentListPage } from "./pages/StudentList";
import { StudentProfilePage } from "./pages/StudentProfile";

function RegistrationRequestsWrapper() {
  const { registrationVersion } = useOutletContext<{ registrationVersion: number }>();
  return <RegistrationRequestsPage registrationVersion={registrationVersion} />;
}

function App() {
  return (
    <Routes>
      <Route element={<Layout />} path="/">
        <Route element={<StudentListPage />} index />
        <Route element={<RegistrationRequestsWrapper />} path="registrations" />
        <Route element={<AddStudentPage />} path="add" />
        <Route element={<StudentProfilePage />} path="students/:id" />
        <Route element={<Navigate to="/" replace />} path="*" />
      </Route>
    </Routes>
  );
}

export default App;
