import { Navigate, Route, Routes } from "react-router-dom";

import { Layout } from "./components/Layout";
import { AddStudentPage } from "./pages/AddStudent";
import { StudentListPage } from "./pages/StudentList";
import { StudentProfilePage } from "./pages/StudentProfile";

function App() {
  return (
    <Routes>
      <Route element={<Layout />} path="/">
        <Route element={<StudentListPage />} index />
        <Route element={<AddStudentPage />} path="add" />
        <Route element={<StudentProfilePage />} path="students/:id" />
        <Route element={<Navigate to="/" replace />} path="*" />
      </Route>
    </Routes>
  );
}

export default App;
