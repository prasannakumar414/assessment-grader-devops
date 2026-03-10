import { useCallback, useEffect, useMemo, useState } from "react";
import { Link } from "react-router-dom";

import { checkDocker, listStudents } from "../api/client";
import { StatusBadge } from "../components/StatusBadge";
import type { Student } from "../types/student";

export function StudentListPage({ stageVersion }: { stageVersion: number }) {
  const [students, setStudents] = useState<Student[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");
  const [checkingDockerIds, setCheckingDockerIds] = useState<Set<number>>(new Set());

  const fetchStudents = useCallback(async () => {
    setLoading(true);
    setError("");
    try {
      const data = await listStudents(true);
      setStudents(data);
    } catch (err) {
      setError("Failed to load students.");
      console.error(err);
    } finally {
      setLoading(false);
    }
  }, []);

  const handleCheckDocker = useCallback(async (id: number) => {
    setCheckingDockerIds((prev) => new Set(prev).add(id));
    try {
      const result = await checkDocker(id);
      if (!result.passed) {
        setError(`Docker check failed for student ${id}: ${result.errorMessage}`);
      }
      await fetchStudents();
    } catch (err) {
      setError("Failed to run Docker check.");
      console.error(err);
    } finally {
      setCheckingDockerIds((prev) => {
        const next = new Set(prev);
        next.delete(id);
        return next;
      });
    }
  }, [fetchStudents]);

  useEffect(() => {
    fetchStudents();
  }, [fetchStudents, stageVersion]);

  const summary = useMemo(() => {
    const total = students.length;
    const github = students.filter((s) => s.githubStatus === "passed").length;
    const docker = students.filter((s) => s.dockerStatus === "passed").length;
    const k8s = students.filter((s) => s.k8sStatus === "passed").length;
    const allDone = students.filter(
      (s) =>
        s.githubStatus === "passed" &&
        s.dockerStatus === "passed" &&
        s.k8sStatus === "passed"
    ).length;
    return { total, github, docker, k8s, allDone };
  }, [students]);

  return (
    <section className="space-y-4">
      <div className="flex flex-wrap items-center justify-between gap-3">
        <div>
          <h2 className="text-xl font-semibold">Students</h2>
          <p className="text-sm text-slate-600">
            Total: {summary.total} | GitHub: {summary.github}/{summary.total} |
            Docker: {summary.docker}/{summary.total} | K8s: {summary.k8s}/{summary.total} |
            All Complete: {summary.allDone}/{summary.total}
          </p>
        </div>
        <div className="flex gap-2">
          <Link className="rounded border border-slate-300 px-3 py-2 text-sm font-medium" to="/add">
            Add Student
          </Link>
          <button
            className="rounded border border-slate-300 px-3 py-2 text-sm font-medium"
            onClick={fetchStudents}
            disabled={loading}
            type="button"
          >
            {loading ? "Refreshing..." : "Refresh"}
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
              <th className="px-3 py-3 text-center">GitHub</th>
              <th className="px-3 py-3 text-center">Docker</th>
              <th className="px-3 py-3 text-center">K8s</th>
              <th className="px-4 py-3">Actions</th>
            </tr>
          </thead>
          <tbody>
            {loading && students.length === 0 ? (
              <tr>
                <td className="px-4 py-6 text-slate-500" colSpan={8}>
                  Loading students...
                </td>
              </tr>
            ) : students.length === 0 ? (
              <tr>
                <td className="px-4 py-6 text-slate-500" colSpan={8}>
                  No approved students yet.
                </td>
              </tr>
            ) : (
              students.map((student) => (
                <tr className="border-t border-slate-100" key={student.id}>
                  <td className="px-4 py-3 font-medium">{student.name}</td>
                  <td className="px-4 py-3">{student.email}</td>
                  <td className="px-4 py-3 font-mono text-xs">{student.githubUsername}</td>
                  <td className="px-4 py-3 font-mono text-xs">{student.dockerHubUsername}</td>
                  <td className="px-3 py-3 text-center">
                    <StatusBadge status={student.githubStatus} />
                  </td>
                  <td className="px-3 py-3 text-center">
                    <StatusBadge status={student.dockerStatus} />
                  </td>
                  <td className="px-3 py-3 text-center">
                    <StatusBadge status={student.k8sStatus} />
                  </td>
                  <td className="px-4 py-3">
                    <div className="flex gap-1">
                      <Link
                        className="rounded border border-slate-300 px-2 py-1 text-xs"
                        to={`/students/${student.id}`}
                      >
                        Profile
                      </Link>
                      <button
                        className="rounded border border-blue-400 bg-blue-50 px-2 py-1 text-xs text-blue-700 hover:bg-blue-100 disabled:opacity-50"
                        onClick={() => handleCheckDocker(student.id)}
                        disabled={checkingDockerIds.has(student.id)}
                        type="button"
                      >
                        {checkingDockerIds.has(student.id) ? "Checking..." : "Check Docker"}
                      </button>
                    </div>
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
