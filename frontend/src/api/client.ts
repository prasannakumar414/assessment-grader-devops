import axios from "axios";

import type {
  CreateStudentPayload,
  RunCheckResult,
  Student,
  StudentStatus,
} from "../types/student";

const api = axios.create({
  baseURL: "/api",
});

export async function createStudent(payload: CreateStudentPayload): Promise<Student> {
  const response = await api.post<Student>("/students", payload);
  return response.data;
}

export async function listStudents(status?: StudentStatus): Promise<Student[]> {
  const response = await api.get<Student[]>("/students", {
    params: status ? { status } : undefined,
  });
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

export async function runCheckAll(): Promise<RunCheckResult[]> {
  const response = await api.post<{ count: number; results: RunCheckResult[] }>(
    "/run-check"
  );
  return response.data.results;
}

export async function runCheckById(id: number): Promise<RunCheckResult> {
  const response = await api.post<RunCheckResult>(`/run-check/${id}`);
  return response.data;
}
