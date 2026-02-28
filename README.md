# Docker Workshop Assessment Grader

Full-stack application to assess student Docker submissions.

## Tech Stack

- Backend: Go + Gin + GORM + SQLite
- Frontend: React + Vite + Tailwind CSS
- Container checks: Docker Engine API (Go SDK)

## Features

- Add and manage students (`name`, `email`, `rollNo`)
- Trigger assessment checks for one or all students
- Pull student Docker image (using roll number naming convention)
- Run container and call `GET /api/info` on port `8080`
- Mark student as `passed` when response email matches registered email
- Track `pending`, `failed`, and `passed` statuses

## API Endpoints

- `POST /api/students`
- `GET /api/students`
- `GET /api/students/:id`
- `PUT /api/students/:id`
- `DELETE /api/students/:id`
- `POST /api/run-check`
- `POST /api/run-check/:id`

## Development

### Prerequisites

- Go 1.22+
- Node.js 20+
- Docker running locally

### 1) Run backend

```bash
go run .
```

Backend runs on `http://localhost:8080`.

### 2) Run frontend (dev server)

```bash
cd frontend
npm install
npm run dev
```

Frontend runs on `http://localhost:5173` and proxies `/api` to backend.

## Production Build (Go serving React)

Build frontend first, then run backend:

```bash
cd frontend
npm install
npm run build
cd ..
go run .
```

Go server embeds and serves files from `frontend/dist`.

