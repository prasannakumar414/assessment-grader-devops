import type { StudentStatus } from "../types/student";

const statusStyles: Record<StudentStatus, string> = {
  pending: "bg-slate-200 text-slate-700",
  passed: "bg-green-100 text-green-800",
  failed: "bg-red-100 text-red-800",
};

export function StatusBadge({ status }: { status: StudentStatus }) {
  return (
    <span
      className={`inline-flex rounded-full px-2 py-0.5 text-xs font-semibold capitalize ${statusStyles[status]}`}
    >
      {status}
    </span>
  );
}
