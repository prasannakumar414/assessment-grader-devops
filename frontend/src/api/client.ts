import axios from "axios";

import { clearToken, getToken } from "../auth";
import type { CreateStudentPayload, Student } from "../types/student";

const api = axios.create({
  baseURL: "/api",
});

api.interceptors.request.use((config) => {
  const token = getToken();
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      clearToken();
      window.location.href = "/login";
    }
    return Promise.reject(error);
  }
);

export async function login(
  username: string,
  password: string
): Promise<string> {
  const response = await api.post<{ token: string }>("/auth/login", {
    username,
    password,
  });
  return response.data.token;
}

export async function createStudent(
  payload: CreateStudentPayload
): Promise<Student> {
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
  const response = await api.post<{ approved: number }>(
    "/registrations/approve-all"
  );
  return response.data;
}

export async function checkDocker(
  id: number
): Promise<{ passed: boolean; errorMessage: string }> {
  const response = await api.post<{ passed: boolean; errorMessage: string }>(
    `/students/${id}/check-docker`
  );
  return response.data;
}
