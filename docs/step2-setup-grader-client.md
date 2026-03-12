# Step 2 - Set Up Grader Client and Register

## Goal

Run the grader client on your machine, open the UI, and register your details.

## Option A - Pull Pre-Built Grader Client Image

```bash
docker pull prasannakumar08/docker-assesment-test-notify
docker run -d \
  -p 3000:3000 \
  -v /var/run/docker.sock:/var/run/docker.sock \
  -e SERVER_URL=http://<instructor-server-ip>:8080 \
  -e IMAGE_NAME=docker-assessment-test \
  -e GITHUB_REPO=docker-assessment-test \
  prasannakumar08/docker-assesment-test-notify
```

## Option B - Build Grader Client Locally

From this repository root:

```bash
docker build -f Dockerfile.client -t grader-client .
docker run -d \
  -p 3000:3000 \
  -v /var/run/docker.sock:/var/run/docker.sock \
  -e SERVER_URL=http://<instructor-server-ip>:8080 \
  -e IMAGE_NAME=docker-assessment-test \
  -e GITHUB_REPO=docker-assessment-test \
  grader-client
```

## Why the Docker Socket Mount Is Required

`-v /var/run/docker.sock:/var/run/docker.sock` allows the grader client to:

- pull your Docker Hub image
- run your container
- verify your `/api/info` endpoint for Stage 2

Without this mount, Docker check will fail.

## Open the UI and Register

1. Open `http://localhost:3000`
2. Fill the form:
   - Student Name
   - GitHub Username
   - Docker Hub Username
   - GitHub Repo Name
   - Email
3. Click **Register**

After registration:

- If not yet approved: you see "Waiting for admin approval"
- Click **Register** again later to refresh status
- Once approved, **Check GitHub** and **Check Docker** buttons are enabled

## Important Matching Rules

- `Email` must be the same email your app returns in checks
- `GitHub Username` and `GitHub Repo Name` must match your public repository
- `Docker Hub Username` must match where you push your image
- `IMAGE_NAME` env var must match the repository name you push on Docker Hub

Continue to [Step 3 - Run GitHub Check](./step3-check-github.md).
