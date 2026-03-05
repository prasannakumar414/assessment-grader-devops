import { useCallback, useEffect, useState } from "react";

import { approveAll, approveStudent, listStudents } from "../api/client";
import type { Student } from "../types/student";

export function RegistrationRequestsPage({ registrationVersion }: { registrationVersion: number }) {
  const [students, setStudents] = useState<Student[]>([]);
  const [loading, setLoading] = useState(false);
  const [approvingId, setApprovingId] = useState<number | null>(null);
  const [approvingAll, setApprovingAll] = useState(false);
  const [error, setError] = useState("");

  const fetchPending = useCallback(async () => {
    setLoading(true);
    setError("");
    try {
      const data = await listStudents(false);
      setStudents(data);
    } catch (err) {
      setError("Failed to load registration requests.");
      console.error(err);
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    fetchPending();
  }, [fetchPending, registrationVersion]);

  async function onApprove(id: number) {
    setApprovingId(id);
    setError("");
    try {
      await approveStudent(id);
      setStudents((prev) => prev.filter((s) => s.id !== id));
    } catch (err) {
      setError("Failed to approve student.");
      console.error(err);
    } finally {
      setApprovingId(null);
    }
  }

  async function onApproveAll() {
    setApprovingAll(true);
    setError("");
    try {
      await approveAll();
      setStudents([]);
    } catch (err) {
      setError("Failed to approve all students.");
      console.error(err);
    } finally {
      setApprovingAll(false);
    }
  }

  return (
    <section className="space-y-4">
      <div className="flex flex-wrap items-center justify-between gap-3">
        <div>
          <h2 className="text-xl font-semibold">Registration Requests</h2>
          <p className="text-sm text-slate-600">{students.length} pending approval</p>
        </div>
        <div className="flex gap-2">
          <button
            className="rounded border border-slate-300 px-3 py-2 text-sm font-medium"
            onClick={fetchPending}
            disabled={loading}
            type="button"
          >
            {loading ? "Refreshing..." : "Refresh"}
          </button>
          <button
            className="rounded bg-green-600 px-3 py-2 text-sm font-medium text-white disabled:opacity-60"
            disabled={approvingAll || students.length === 0}
            onClick={onApproveAll}
            type="button"
          >
            {approvingAll ? "Approving..." : "Approve All"}
          </button>
        </div>
      </div>

      {error ? <p className="rounded bg-red-100 p-3 text-sm text-red-700">{error}</p> : null}

      <div className="overflow-x-auto rounded-lg border border-slate-200 bg-white shadow-sm">
        <table className="min-w-full text-left text-sm">
          <thead className="bg-slate-100 text-slate-700">
            <tr>
              <th className="px-4 py-3">Name</th>
              <th className="px-4 py-3">Email</th>
              <th className="px-4 py-3">GitHub</th>
              <th className="px-4 py-3">Docker Hub</th>
              <th className="px-4 py-3">Registered</th>
              <th className="px-4 py-3">Action</th>
            </tr>
          </thead>
          <tbody>
            {loading && students.length === 0 ? (
              <tr>
                <td className="px-4 py-6 text-slate-500" colSpan={6}>
                  Loading...
                </td>
              </tr>
            ) : students.length === 0 ? (
              <tr>
                <td className="px-4 py-6 text-slate-500" colSpan={6}>
                  No pending registration requests.
                </td>
              </tr>
            ) : (
              students.map((student) => (
                <tr className="border-t border-slate-100" key={student.id}>
                  <td className="px-4 py-3 font-medium">{student.name}</td>
                  <td className="px-4 py-3">{student.email}</td>
                  <td className="px-4 py-3 font-mono text-xs">{student.githubUsername}</td>
                  <td className="px-4 py-3 font-mono text-xs">{student.dockerHubUsername}</td>
                  <td className="px-4 py-3 text-xs text-slate-500">
                    {new Date(student.createdAt).toLocaleString()}
                  </td>
                  <td className="px-4 py-3">
                    <button
                      className="rounded bg-green-600 px-3 py-1 text-xs font-medium text-white disabled:opacity-60"
                      disabled={approvingId === student.id}
                      onClick={() => onApprove(student.id)}
                      type="button"
                    >
                      {approvingId === student.id ? "Approving..." : "Approve"}
                    </button>
                  </td>
                </tr>
              ))
            )}
          </tbody>
        </table>
      </div>
    </section>
  );
}
