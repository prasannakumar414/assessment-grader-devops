import { useEffect, useState } from "react";
import { Link, useParams } from "react-router-dom";

import { getStudent, runCheckById } from "../api/client";
import { StatusBadge } from "../components/StatusBadge";
import type { Student } from "../types/student";

export function StudentProfilePage() {
  const { id } = useParams();
  const [student, setStudent] = useState<Student | null>(null);
  const [loading, setLoading] = useState(true);
  const [running, setRunning] = useState(false);
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

  async function onRunCheck() {
    if (!student) return;
    setRunning(true);
    setError("");
    try {
      await runCheckById(student.id);
      await fetchStudent();
    } catch (err) {
      setError("Failed to run check.");
      console.error(err);
    } finally {
      setRunning(false);
    }
  }

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

  return (
    <section className="space-y-4">
      <div className="flex items-center justify-between">
        <h2 className="text-xl font-semibold">{student.name}</h2>
        <Link className="rounded border border-slate-300 px-3 py-2 text-sm" to="/">
          Back
        </Link>
      </div>

      {error ? <p className="rounded bg-red-100 p-3 text-sm text-red-700">{error}</p> : null}

      <div className="rounded-lg border border-slate-200 bg-white p-5 shadow-sm">
        <dl className="grid gap-4 sm:grid-cols-2">
          <div>
            <dt className="text-xs uppercase tracking-wide text-slate-500">Email</dt>
            <dd className="mt-1">{student.email}</dd>
          </div>
          <div>
            <dt className="text-xs uppercase tracking-wide text-slate-500">Roll No / Image</dt>
            <dd className="mt-1 font-mono text-xs">{student.rollNo}</dd>
          </div>
          <div>
            <dt className="text-xs uppercase tracking-wide text-slate-500">Status</dt>
            <dd className="mt-1">
              <StatusBadge status={student.status} />
            </dd>
          </div>
          <div>
            <dt className="text-xs uppercase tracking-wide text-slate-500">Last Checked</dt>
            <dd className="mt-1">
              {student.lastCheckedAt ? new Date(student.lastCheckedAt).toLocaleString() : "Not checked yet"}
            </dd>
          </div>
        </dl>

        {student.errorMessage ? (
          <div className="mt-4 rounded bg-red-50 p-3 text-sm text-red-700">
            <strong>Last error:</strong> {student.errorMessage}
          </div>
        ) : null}

        <div className="mt-5">
          <button
            className="rounded bg-slate-900 px-4 py-2 text-sm font-medium text-white disabled:opacity-60"
            disabled={running}
            onClick={onRunCheck}
            type="button"
          >
            {running ? "Running check..." : "Run Check"}
          </button>
        </div>
      </div>
    </section>
  );
}
