import { type FormEvent, useState } from "react";
import { Link, useNavigate } from "react-router-dom";

import { createStudent } from "../api/client";

export function AddStudentPage() {
  const navigate = useNavigate();
  const [name, setName] = useState("");
  const [email, setEmail] = useState("");
  const [githubUsername, setGithubUsername] = useState("");
  const [dockerHubUsername, setDockerHubUsername] = useState("");
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");

  async function onSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setLoading(true);
    setError("");
    try {
      await createStudent({ name, email, githubUsername, dockerHubUsername });
      navigate("/");
    } catch (err) {
      setError("Failed to add student. Please check inputs and try again.");
      console.error(err);
    } finally {
      setLoading(false);
    }
  }

  return (
    <section className="mx-auto max-w-xl rounded-lg border border-slate-200 bg-white p-6 shadow-sm">
      <div className="mb-4">
        <h2 className="text-xl font-semibold">Add Student</h2>
        <p className="mt-1 text-sm text-slate-600">
          Admin-created students are auto-approved.
        </p>
      </div>

      <form className="space-y-4" onSubmit={onSubmit}>
        <label className="block">
          <span className="mb-1 block text-sm font-medium">Name</span>
          <input
            className="w-full rounded border border-slate-300 px-3 py-2"
            value={name}
            onChange={(event) => setName(event.target.value)}
            placeholder="Student name"
            required
          />
        </label>

        <label className="block">
          <span className="mb-1 block text-sm font-medium">Email</span>
          <input
            className="w-full rounded border border-slate-300 px-3 py-2"
            value={email}
            onChange={(event) => setEmail(event.target.value)}
            placeholder="student@example.com"
            type="email"
            required
          />
        </label>

        <label className="block">
          <span className="mb-1 block text-sm font-medium">GitHub Username</span>
          <input
            className="w-full rounded border border-slate-300 px-3 py-2"
            value={githubUsername}
            onChange={(event) => setGithubUsername(event.target.value)}
            placeholder="johndoe"
            required
          />
        </label>

        <label className="block">
          <span className="mb-1 block text-sm font-medium">Docker Hub Username</span>
          <input
            className="w-full rounded border border-slate-300 px-3 py-2"
            value={dockerHubUsername}
            onChange={(event) => setDockerHubUsername(event.target.value)}
            placeholder="johndoe"
            required
          />
        </label>

        {error ? <p className="text-sm text-red-600">{error}</p> : null}

        <div className="flex gap-3">
          <button
            disabled={loading}
            className="rounded bg-slate-900 px-4 py-2 text-sm font-medium text-white disabled:opacity-60"
            type="submit"
          >
            {loading ? "Saving..." : "Add Student"}
          </button>
          <Link className="rounded border border-slate-300 px-4 py-2 text-sm font-medium" to="/">
            Back to list
          </Link>
        </div>
      </form>
    </section>
  );
}
