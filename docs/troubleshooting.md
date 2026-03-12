# Troubleshooting and FAQ

## 1) "Cannot connect to server"

Symptoms:

- Register/check actions fail with connection errors

Fixes:

- Verify grader client container is running: `docker ps`
- Confirm `SERVER_URL` points to instructor server: `http://<instructor-server-ip>:8080`
- Confirm server is reachable from your machine:

```bash
curl http://<instructor-server-ip>:8080/api/register
```

## 2) "Waiting for admin approval"

Symptoms:

- Buttons are disabled
- Status says waiting for approval

Fixes:

- Ask instructor/admin to approve your registration
- Click **Register** again to refresh status

## 3) GitHub check failed

Symptoms:

- "file main.go not found" or failed GitHub stage

Fixes:

- Repository must be public
- File must be exactly `main.go` in root
- Confirm GitHub username and repo name in UI are correct
- Verify API directly:

```bash
curl -i https://api.github.com/repos/<your-github-username>/<your-repo>/contents/main.go
```

## 4) Docker check failed

Symptoms:

- Pull errors
- Timeout
- Email mismatch

Fixes:

- Image must be public on Docker Hub
- Image name must match `<dockerHubUsername>/<IMAGE_NAME>:latest`
- App must listen on port `8080`
- App must expose `GET /api/info`
- Response must include your registered email, for example:

```json
{"email":"your-email@example.com"}
```

- Test your image locally:

```bash
docker run --rm -p 8080:8080 -e EMAIL=your-email@example.com <your-dockerhub-username>/docker-assessment-test:latest
curl http://localhost:8080/api/info
```

## 5) Docker permissions issue on Linux

Symptoms:

- Docker commands fail without sudo

Fixes:

- Use `sudo` for Docker commands, or
- Add your user to docker group and re-login

## 6) Kubernetes check not updating

Symptoms:

- K8s stage remains pending/failed

Fixes:

- Verify operator resources are installed:
  - `kubectl get crd apphealths.example.com`
  - `kubectl get deploy -n system`
- Verify `AppHealth` exists:
  - `kubectl get apphealth -o wide`
  - `kubectl describe apphealth my-health-check`
- Verify `reportURL` points to reachable grader server
- Check operator logs:

```bash
kubectl logs -n system deploy/health-operator
```

## 7) `/info` vs `/api/info` path issue

Symptoms:

- App works in Docker stage but fails in Kubernetes stage

Reason:

- Docker stage calls `/api/info`
- Operator default uses `/info`

Fixes:

- Set `infoPath: /api/info` in your `AppHealth` resource, or
- Add `/info` endpoint to your app

## 8) How to reset your registration state in browser

1. Open browser DevTools
2. Go to **Application** (or **Storage**)
3. Open **Local Storage**
4. Remove key: `grader_student`
5. Refresh page and register again

If you still get blocked, share:

- your GitHub repo URL
- your Docker Hub image URL
- exact error message from the UI
- output of relevant `kubectl` or `docker` commands
