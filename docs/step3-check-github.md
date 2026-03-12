# Step 3 - Run GitHub Check

## Goal

Pass the GitHub stage by confirming `main.go` exists in the root of your public repository.

## What the Checker Verifies

The grader calls:

`https://api.github.com/repos/<your-github-username>/<your-repo>/contents/main.go`

If that file exists, GitHub stage passes.

## Steps

1. Confirm your repository is public
2. Confirm `main.go` is in repository root
3. Open grader UI at `http://localhost:3000`
4. Click **Check GitHub**

If successful, you see a green **Passed!** result.

## Common Failures

- Repository is private
- `main.go` is missing
- `main.go` is inside a subfolder instead of root
- GitHub username typo in registration form
- Repo name typo in registration form

## Quick Verification Commands

```bash
# In your app repository
ls
# You should see: main.go
```

```bash
# Optional: test GitHub API directly
curl -i https://api.github.com/repos/<your-github-username>/<your-repo>/contents/main.go
```

Continue to [Step 4 - Build and Push Docker Image, Run Docker Check](./step4-build-and-push-docker.md).
