import axios from "axios";

import type { CreateStudentPayload, Student } from "../types/student";

const api = axios.create({
  baseURL: "/api",
});

export async function createStudent(payload: CreateStudentPayload): Promise<Student> {
  const response = await api.post<Student>("/students", payload);
  return response.data;
}

export async function listStudents(approved?: boolean): Promise<Student[]> {
  const params: Record<string, string> = {};
  if (approved !== undefined) {
    params.approved = String(approved);
  }
  const response = await api.get<Student[]>("/students", { params });
  return response.data;
}

export async function getStudent(id: number): Promise<Student> {
  const response = await api.get<Student>(`/students/${id}`);
  return response.data;
}

export async function updateStudent(
  id: number,
  payload: CreateStudentPayload
): Promise<Student> {
  const response = await api.put<Student>(`/students/${id}`, payload);
  return response.data;
}

export async function deleteStudent(id: number): Promise<void> {
  await api.delete(`/students/${id}`);
}

export async function approveStudent(id: number): Promise<Student> {
  const response = await api.post<Student>(`/registrations/${id}/approve`);
  return response.data;
}

export async function approveAll(): Promise<{ approved: number }> {
  const response = await api.post<{ approved: number }>("/registrations/approve-all");
  return response.data;
}
