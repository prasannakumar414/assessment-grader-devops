# Step 4 - Build and Push Docker Image, Run Docker Check

## Goal

Build a Docker image from your cloned app, push it to Docker Hub, then pass Docker validation in the grader UI.

## Dockerfile (Starter Repo)

The starter repo already includes a working Dockerfile:

```dockerfile
FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY go.mod ./
COPY main.go ./
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/server .

FROM alpine:3.20
WORKDIR /app
COPY --from=builder /out/server /app/server
EXPOSE 8080
ENTRYPOINT ["/app/server"]
```

## Build, Login, Push

From your cloned app repo:

```bash
docker build -t <your-dockerhub-username>/docker-assessment-test:latest .
docker login
docker push <your-dockerhub-username>/docker-assessment-test:latest
```

## Image Naming Must Match

The grader checks this exact image format:

`<dockerHubUsername>/<IMAGE_NAME>:latest`

Where:

- `<dockerHubUsername>` comes from the form
- `<IMAGE_NAME>` comes from grader client env var (example: `docker-assessment-test`)

Example resolved image:

`johndoe/docker-assessment-test:latest`

## Run Docker Check

1. Open `http://localhost:3000`
2. Click **Check Docker**

## What the System Verifies

- Pulls your image from Docker Hub
- Starts a container
- Calls `GET /api/info` on port `8080`
- Validates returned email equals your registered email

Expected response format:

```json
{"email":"your-email@example.com"}
```

## Common Failures

- Image is not public
- Wrong image name (`IMAGE_NAME` mismatch)
- App is not listening on port `8080`
- `/api/info` missing or returns invalid JSON
- Email does not match registered email

## Quick Local Test Before Clicking Check

```bash
docker run --rm -p 8080:8080 -e EMAIL=your-email@example.com <your-dockerhub-username>/docker-assessment-test:latest
```

In another terminal:

```bash
curl http://localhost:8080/api/info
```

Continue to [Step 5 - Deploy to Kubernetes and Complete K8s Check](./step5-kubernetes.md).
