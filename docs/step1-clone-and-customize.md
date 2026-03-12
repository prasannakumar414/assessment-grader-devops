# Step 1 - Clone and Customize

## Goal

Clone the starter repository, replace the email with your own, and push it to your own public GitHub repository.

Starter repository: [prasannakumar414/docker-assessment-test](https://github.com/prasannakumar414/docker-assessment-test)

## 1) Clone the Starter Repo

```bash
git clone https://github.com/prasannakumar414/docker-assessment-test.git
cd docker-assessment-test
```

## 2) Update Email in the App

The starter app serves `GET /api/info` and returns JSON like:

```json
{"email":"your-email@example.com"}
```

Open `main.go` and update the default email value to your own email.

Current behavior in starter app:

- Reads `EMAIL` environment variable
- Falls back to default `test@example.com` if not set
- Exposes `GET /api/info` on port `8080`

## 3) Commit Your Changes

```bash
git add main.go
git commit -m "Update email for workshop"
```

## 4) Push to Your Own Public GitHub Repo

Create a public repository in your own GitHub account (example: `docker-assessment-test`) and push:

```bash
git remote rename origin upstream
git remote add origin https://github.com/<your-github-username>/docker-assessment-test.git
git push -u origin master
```

## 5) Confirm Before Moving On

- Repo is public
- `main.go` is in repository root
- Your GitHub username and repo name are noted for registration

Continue to [Step 2 - Set Up Grader Client and Register](./step2-setup-grader-client.md).
