# DevOps Assessment Grader

A two-component system for running and grading multi-stage DevOps workshop assessments. An **admin server** manages students and tracks progress, while a **client app** (distributed as a Docker image) runs on each student's machine to execute checks and report results.

## Architecture

```
┌──────────────────────────────┐     ┌──────────────────────────────┐
│      Student Machine         │     │        Admin Server          │
│                              │     │                              │
│  ┌────────────────────────┐  │     │  ┌────────────────────────┐  │
│  │  Client Container      │  │     │  │  Go Server (main.go)   │  │
│  │  (grader-client)       │──┼────>│  │  - Student CRUD        │  │
│  │                        │  │     │  │  - POST /api/register   │  │
│  │  HTML form + checks    │  │     │  │  - POST /api/notify     │  │
│  │  - GitHub API check    │  │     │  │  - GET /api/events (SSE)│  │
│  │  - Docker check        │  │     │  │  - Approval endpoints   │  │
│  └──────────┬─────────────┘  │     │  └────────────┬───────────┘  │
│             │                │     │               │              │
│  ┌──────────▼─────────────┐  │     │  ┌────────────▼───────────┐  │
│  │  Host Docker Engine    │  │     │  │  SQLite + React UI     │  │
│  │  (via socket mount)    │  │     │  │  (embedded frontend)   │  │
│  └────────────────────────┘  │     │  └────────────────────────┘  │
└──────────────────────────────┘     └──────────────────────────────┘
                                                   ▲
                                     ┌─────────────┘
                                     │  K8s Operator
                                     │  POST /api/notify stage=k8s
                                     └─────────────────────────────
```

## Three Assessment Stages

| Stage | Where it runs | What it checks |
|-------|--------------|----------------|
| **GitHub** | Client | Checks if `main.go` exists in `github.com/<username>/<repo>` via GitHub API |
| **Docker** | Client | Pulls `<dockerhub_username>/<image_name>`, runs it, verifies `/api/info` returns the student's email |
| **Kubernetes** | K8s Operator | External operator calls the notify API directly |

Each stage is tracked independently per student. The admin dashboard shows per-stage pass/fail status.

## Registration & Approval Flow

1. Student fills form on the client UI and clicks **Register**
2. Server creates the student with `approved: false`
3. Admin sees the request on the **Registration Requests** page (updated in real time via SSE)
4. Admin clicks **Approve** (or **Approve All**)
5. Student can now run checks — the client verifies approval before each check

## Real-Time Celebrations

When a student passes a stage, the admin dashboard shows a celebration modal:
- **Single stage**: slide-up notification with student name and stage
- **All stages complete**: full-screen confetti celebration

Popup dismiss times are configurable in `frontend/src/config.ts`. Multiple events queue up and display sequentially.

---

## Server (Admin Dashboard)

### Tech Stack

- Go + Gin + GORM + SQLite
- React + Vite + Tailwind CSS (embedded in the Go binary)
- Server-Sent Events for real-time updates

### Prerequisites

- Go 1.25+
- Node.js 20+

### API Endpoints

| Method | Path | Description |
|--------|------|-------------|
| `POST` | `/api/students` | Create student (admin, auto-approved) |
| `GET` | `/api/students` | List students (`?approved=true/false`) |
| `GET` | `/api/students/:id` | Get student |
| `PUT` | `/api/students/:id` | Update student |
| `DELETE` | `/api/students/:id` | Delete student |
| `POST` | `/api/register` | Client registration (creates unapproved student) |
| `POST` | `/api/notify` | Report stage result (`stage`, `email`, `passed`, `errorMessage`) |
| `GET` | `/api/events` | SSE stream (`new_registration`, `stage_complete`, `all_complete`) |
| `POST` | `/api/registrations/:id/approve` | Approve one student |
| `POST` | `/api/registrations/approve-all` | Approve all pending students |

### Development

Run the backend:

```bash
go run .
```

Backend runs on `http://localhost:8080`.

Run the frontend dev server (with API proxy):

```bash
cd frontend
npm install
npm run dev
```

Frontend runs on `http://localhost:5173` and proxies `/api` to the backend.

### Production Build

```bash
cd frontend && npm install && npm run build && cd ..
go run .
```

The Go server embeds `frontend/dist` and serves everything from a single binary.

---

## Client (Student Tool)

### Running via Docker (Recommended)

```bash
docker run -d \
  -p 3000:3000 \
  -v /var/run/docker.sock:/var/run/docker.sock \
  -e SERVER_URL=http://10.0.0.1:8080 \
  -e IMAGE_NAME=workshop-app \
  -e GITHUB_REPO=devops-workshop \
  grader-client
```

Then open `http://localhost:3000` in a browser.

### Configuration

| Env Variable | config.json key | Description |
|-------------|-----------------|-------------|
| `SERVER_URL` | `server_url` | Admin server base URL |
| `IMAGE_NAME` | `image_name` | Docker image name suffix (e.g., `workshop-app`) |
| `GITHUB_REPO` | `github_repo` | GitHub repo name to check for `main.go` |
| `PORT` | — | Client UI port (default: `3000`) |
| `CONFIG_PATH` | — | Path to config.json (default: `./config.json`) |

Environment variables take priority. Missing values fall back to `config.json`.

### Building the Docker Image

```bash
docker build -f Dockerfile.client -t grader-client .
```

### Running Locally (Development)

```bash
export SERVER_URL=http://localhost:8080
export IMAGE_NAME=workshop-app
export GITHUB_REPO=devops-workshop
go run ./cmd/client
```

---

## Project Structure

```
├── main.go                         # Server entry point
├── Dockerfile.client               # Client Docker image build
├── go.mod / go.sum
├── internal/
│   ├── database/db.go              # SQLite + GORM setup
│   ├── docker/runner.go            # Docker image pull/run/verify
│   ├── github/checker.go           # GitHub API file check
│   ├── sse/hub.go                  # SSE broadcast hub
│   ├── models/student.go           # Student model (per-stage status)
│   └── handlers/
│       ├── student.go              # Student CRUD
│       ├── register.go             # Client registration
│       ├── approval.go             # Admin approval
│       ├── notify.go               # Stage result notification
│       └── events.go               # SSE endpoint
├── cmd/client/
│   ├── main.go                     # Client entry point
│   └── frontend.html               # Client UI (embedded)
└── frontend/                       # Admin React app
    ├── src/
    │   ├── api/client.ts
    │   ├── config.ts               # Popup dismiss times
    │   ├── hooks/useSSE.ts
    │   ├── components/
    │   │   ├── Layout.tsx
    │   │   ├── StatusBadge.tsx
    │   │   ├── StageCompleteModal.tsx
    │   │   └── AllCompleteModal.tsx
    │   ├── pages/
    │   │   ├── StudentList.tsx
    │   │   ├── AddStudent.tsx
    │   │   ├── StudentProfile.tsx
    │   │   └── RegistrationRequests.tsx
    │   └── types/student.ts
    └── package.json
```
