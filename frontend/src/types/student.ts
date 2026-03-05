export type StudentStatus = "pending" | "passed" | "failed";
export type StageName = "github" | "docker" | "k8s";

export interface Student {
  id: number;
  name: string;
  email: string;
  githubUsername: string;
  dockerHubUsername: string;
  approved: boolean;

  githubStatus: StudentStatus;
  githubErrorMessage: string;
  githubLastCheckedAt: string | null;

  dockerStatus: StudentStatus;
  dockerErrorMessage: string;
  dockerLastCheckedAt: string | null;

  k8sStatus: StudentStatus;
  k8sErrorMessage: string;
  k8sLastCheckedAt: string | null;

  createdAt: string;
  updatedAt: string;
}

export interface CreateStudentPayload {
  name: string;
  email: string;
  githubUsername: string;
  dockerHubUsername: string;
}

export interface StageCompleteEvent {
  studentName: string;
  stageName: StageName;
}

export interface AllCompleteEvent {
  studentName: string;
}

export interface NewRegistrationEvent {
  studentName: string;
  studentId: number;
}
