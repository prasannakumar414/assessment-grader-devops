import { useCallback, useEffect, useMemo, useState } from "react";
import { Link } from "react-router-dom";

import { listStudents, runCheckAll, runCheckById } from "../api/client";
import { StatusBadge } from "../components/StatusBadge";
import type { Student, StudentStatus } from "../types/student";

type FilterValue = "all" | StudentStatus;

const filters: FilterValue[] = ["all", "pending", "passed", "failed"];

export function StudentListPage() {
  const [students, setStudents] = useState<Student[]>([]);
  const [filter, setFilter] = useState<FilterValue>("all");
  const [loading, setLoading] = useState(false);
  const [runningAll, setRunningAll] = useState(false);
  const [activeRow, setActiveRow] = useState<number | null>(null);
  const [error, setError] = useState("");

  const fetchStudents = useCallback(async () => {
    setLoading(true);
    setError("");
    try {
      const data = await listStudents(filter === "all" ? undefined : filter);
      setStudents(data);
    } catch (err) {
      setError("Failed to load students.");
      console.error(err);
    } finally {
      setLoading(false);
    }
  }, [filter]);

  useEffect(() => {
    fetchStudents();
  }, [fetchStudents]);

  const summary = useMemo(() => {
    const result = { total: 0, pending: 0, passed: 0, failed: 0 };
    students.forEach((student) => {
      result.total += 1;
      result[student.status] += 1;
    });
    return result;
  }, [students]);

  async function onRunCheckAll() {
    setRunningAll(true);
    setError("");
    try {
      await runCheckAll();
      await fetchStudents();
    } catch (err) {
      setError("Failed to run checks for all students.");
      console.error(err);
    } finally {
      setRunningAll(false);
    }
  }

  async function onRunCheckOne(id: number) {
    setActiveRow(id);
    setError("");
    try {
      await runCheckById(id);
      await fetchStudents();
    } catch (err) {
      setError("Failed to run check for selected student.");
      console.error(err);
    } finally {
      setActiveRow(null);
    }
  }

  return (
    <section className="space-y-4">
      <div className="flex flex-wrap items-center justify-between gap-3">
        <div>
          <h2 className="text-xl font-semibold">Students</h2>
          <p className="text-sm text-slate-600">
            Total: {summary.total} | Passed: {summary.passed} | Failed: {summary.failed} | Pending:{" "}
            {summary.pending}
          </p>
        </div>
        <div className="flex gap-2">
          <Link className="rounded border border-slate-300 px-3 py-2 text-sm font-medium" to="/add">
            Add Student
          </Link>
          <button
            className="rounded bg-slate-900 px-3 py-2 text-sm font-medium text-white disabled:opacity-60"
            disabled={runningAll}
            onClick={onRunCheckAll}
          >
            {runningAll ? "Running checks..." : "Run Check All"}
          </button>
        </div>
      </div>

      <div className="flex flex-wrap gap-2">
        {filters.map((value) => (
          <button
            key={value}
            className={`rounded-full px-3 py-1 text-sm capitalize ${
              filter === value ? "bg-slate-900 text-white" : "bg-white text-slate-700"
            }`}
            onClick={() => setFilter(value)}
            type="button"
          >
            {value}
          </button>
        ))}
      </div>

      {error ? <p className="rounded bg-red-100 p-3 text-sm text-red-700">{error}</p> : null}

      <div className="overflow-x-auto rounded-lg border border-slate-200 bg-white shadow-sm">
        <table className="min-w-full text-left text-sm">
          <thead className="bg-slate-100 text-slate-700">
            <tr>
              <th className="px-4 py-3">Name</th>
              <th className="px-4 py-3">Roll No</th>
              <th className="px-4 py-3">Email</th>
              <th className="px-4 py-3">Status</th>
              <th className="px-4 py-3">Actions</th>
            </tr>
          </thead>
          <tbody>
            {loading ? (
              <tr>
                <td className="px-4 py-6 text-slate-500" colSpan={5}>
                  Loading students...
                </td>
              </tr>
            ) : students.length === 0 ? (
              <tr>
                <td className="px-4 py-6 text-slate-500" colSpan={5}>
                  No students found.
                </td>
              </tr>
            ) : (
              students.map((student) => (
                <tr className="border-t border-slate-100" key={student.id}>
                  <td className="px-4 py-3">{student.name}</td>
                  <td className="px-4 py-3 font-mono text-xs">{student.rollNo}</td>
                  <td className="px-4 py-3">{student.email}</td>
                  <td className="px-4 py-3">
                    <StatusBadge status={student.status} />
                  </td>
                  <td className="px-4 py-3">
                    <div className="flex flex-wrap gap-2">
                      <Link className="rounded border border-slate-300 px-2 py-1 text-xs" to={`/students/${student.id}`}>
                        Profile
                      </Link>
                      <button
                        className="rounded border border-slate-300 px-2 py-1 text-xs disabled:opacity-60"
                        disabled={activeRow === student.id}
                        onClick={() => onRunCheckOne(student.id)}
                        type="button"
                      >
                        {activeRow === student.id ? "Checking..." : "Re-check"}
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
