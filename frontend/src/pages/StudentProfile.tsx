import { useEffect, useState } from "react";
import { Link, useParams } from "react-router-dom";

import { getStudent } from "../api/client";
import { StatusBadge } from "../components/StatusBadge";
import type { Student, StageName, StudentStatus } from "../types/student";

const stages: { key: StageName; label: string }[] = [
  { key: "github", label: "GitHub" },
  { key: "docker", label: "Docker" },
  { key: "k8s", label: "Kubernetes" },
];

function stageInfo(student: Student, stage: StageName) {
  switch (stage) {
    case "github":
      return { status: student.githubStatus, error: student.githubErrorMessage, checkedAt: student.githubLastCheckedAt };
    case "docker":
      return { status: student.dockerStatus, error: student.dockerErrorMessage, checkedAt: student.dockerLastCheckedAt };
    case "k8s":
      return { status: student.k8sStatus, error: student.k8sErrorMessage, checkedAt: student.k8sLastCheckedAt };
  }
}

export function StudentProfilePage() {
  const { id } = useParams();
  const [student, setStudent] = useState<Student | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");

  async function fetchStudent() {
    if (!id) return;
    setLoading(true);
    setError("");
    try {
      const data = await getStudent(Number(id));
      setStudent(data);
    } catch (err) {
      setError("Failed to load student profile.");
      console.error(err);
    } finally {
      setLoading(false);
    }
  }

  useEffect(() => {
    fetchStudent();
  }, [id]);

  if (loading) {
    return <p>Loading profile...</p>;
  }

  if (!student) {
    return (
      <section className="space-y-3">
        <p className="text-red-700">{error || "Student not found."}</p>
        <Link className="text-sm text-slate-700 underline" to="/">
          Back to students
        </Link>
      </section>
    );
  }

  const passedCount = stages.filter((s) => stageInfo(student, s.key).status === "passed").length;

  return (
    <section className="space-y-4">
      <div className="flex items-center justify-between">
        <h2 className="text-xl font-semibold">{student.name}</h2>
        <div className="flex gap-2">
          <button
            className="rounded border border-slate-300 px-3 py-2 text-sm"
            onClick={fetchStudent}
            type="button"
          >
            Refresh
          </button>
          <Link className="rounded border border-slate-300 px-3 py-2 text-sm" to="/">
            Back
          </Link>
        </div>
      </div>

      {error ? <p className="rounded bg-red-100 p-3 text-sm text-red-700">{error}</p> : null}

      <div className="rounded-lg border border-slate-200 bg-white p-5 shadow-sm">
        <dl className="grid gap-4 sm:grid-cols-2">
          <div>
            <dt className="text-xs uppercase tracking-wide text-slate-500">Email</dt>
            <dd className="mt-1">{student.email}</dd>
          </div>
          <div>
            <dt className="text-xs uppercase tracking-wide text-slate-500">GitHub Username</dt>
            <dd className="mt-1 font-mono text-sm">{student.githubUsername}</dd>
          </div>
          <div>
            <dt className="text-xs uppercase tracking-wide text-slate-500">Docker Hub Username</dt>
            <dd className="mt-1 font-mono text-sm">{student.dockerHubUsername}</dd>
          </div>
          <div>
            <dt className="text-xs uppercase tracking-wide text-slate-500">Approved</dt>
            <dd className="mt-1">
              <span className={`inline-flex rounded-full px-2 py-0.5 text-xs font-semibold ${student.approved ? "bg-green-100 text-green-800" : "bg-yellow-100 text-yellow-800"}`}>
                {student.approved ? "Yes" : "Pending"}
              </span>
            </dd>
          </div>
        </dl>
      </div>

      <div className="rounded-lg border border-slate-200 bg-white p-5 shadow-sm">
        <h3 className="mb-4 text-sm font-semibold uppercase tracking-wide text-slate-500">
          Stage Progress ({passedCount}/{stages.length})
        </h3>

        <div className="flex items-center gap-2 mb-6">
          {stages.map((s, i) => {
            const info = stageInfo(student, s.key);
            const done = info.status === "passed";
            return (
              <div key={s.key} className="flex items-center gap-2">
                <div className={`flex h-8 w-8 items-center justify-center rounded-full text-xs font-bold ${done ? "bg-green-500 text-white" : "bg-slate-200 text-slate-500"}`}>
                  {done ? "\u2713" : i + 1}
                </div>
                <span className={`text-sm font-medium ${done ? "text-green-700" : "text-slate-500"}`}>{s.label}</span>
                {i < stages.length - 1 && <div className={`h-0.5 w-8 ${done ? "bg-green-400" : "bg-slate-200"}`} />}
              </div>
            );
          })}
        </div>

        <div className="space-y-3">
          {stages.map((s) => {
            const info = stageInfo(student, s.key);
            return (
              <div key={s.key} className="rounded border border-slate-100 p-3">
                <div className="flex items-center justify-between">
                  <span className="text-sm font-medium">{s.label}</span>
                  <StatusBadge status={info.status as StudentStatus} />
                </div>
                <div className="mt-1 text-xs text-slate-500">
                  {info.checkedAt
                    ? `Last checked: ${new Date(info.checkedAt).toLocaleString()}`
                    : "Not checked yet"}
                </div>
                {info.error ? (
                  <div className="mt-1 text-xs text-red-600">{info.error}</div>
                ) : null}
              </div>
            );
          })}
        </div>
      </div>
    </section>
  );
}
