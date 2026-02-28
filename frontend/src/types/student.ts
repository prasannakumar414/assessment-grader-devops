export type StudentStatus = "pending" | "passed" | "failed";

export interface Student {
  id: number;
  name: string;
  email: string;
  rollNo: string;
  status: StudentStatus;
  lastCheckedAt: string | null;
  errorMessage: string;
  createdAt: string;
  updatedAt: string;
}

export interface CreateStudentPayload {
  name: string;
  email: string;
  rollNo: string;
}

export interface RunCheckResult {
  studentId: number;
  rollNo: string;
  status: StudentStatus;
  errorMessage?: string;
}
