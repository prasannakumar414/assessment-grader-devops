package main

import (
	"bytes"
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"docker-workshop-assesment-grader/internal/docker"
	"docker-workshop-assesment-grader/internal/github"
)

//go:embed frontend.html
var frontendHTML embed.FS

type Config struct {
	ServerURL  string `json:"server_url"`
	ImageName  string `json:"image_name"`
	GithubRepo string `json:"github_repo"`
}

type checkRequest struct {
	Name              string `json:"name"`
	Email             string `json:"email"`
	GithubUsername    string `json:"githubUsername"`
	DockerHubUsername string `json:"dockerHubUsername"`
}

func main() {
	cfg := loadConfig()
	log.Printf("client starting — server=%s image=%s repo=%s", cfg.ServerURL, cfg.ImageName, cfg.GithubRepo)

	runner, err := docker.NewRunner()
	if err != nil {
		log.Fatalf("docker runner init failed: %v", err)
	}

	mux := http.NewServeMux()

	mux.HandleFunc("POST /api/register", func(w http.ResponseWriter, r *http.Request) {
		handleRegister(w, r, cfg)
	})
	mux.HandleFunc("POST /api/check/github", func(w http.ResponseWriter, r *http.Request) {
		handleCheckGitHub(w, r, cfg)
	})
	mux.HandleFunc("POST /api/check/docker", func(w http.ResponseWriter, r *http.Request) {
		handleCheckDocker(w, r, cfg, runner)
	})
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		data, _ := frontendHTML.ReadFile("frontend.html")
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write(data)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Printf("client UI at http://localhost:%s", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}

func loadConfig() Config {
	cfg := Config{
		ServerURL:  os.Getenv("SERVER_URL"),
		ImageName:  os.Getenv("IMAGE_NAME"),
		GithubRepo: os.Getenv("GITHUB_REPO"),
	}

	if cfg.ServerURL == "" || cfg.ImageName == "" || cfg.GithubRepo == "" {
		configPath := os.Getenv("CONFIG_PATH")
		if configPath == "" {
			configPath = "config.json"
		}
		if data, err := os.ReadFile(configPath); err == nil {
			var fileCfg Config
			if jsonErr := json.Unmarshal(data, &fileCfg); jsonErr == nil {
				if cfg.ServerURL == "" {
					cfg.ServerURL = fileCfg.ServerURL
				}
				if cfg.ImageName == "" {
					cfg.ImageName = fileCfg.ImageName
				}
				if cfg.GithubRepo == "" {
					cfg.GithubRepo = fileCfg.GithubRepo
				}
			}
		}
	}

	if cfg.ServerURL == "" {
		log.Fatal("SERVER_URL is required (env var or config.json)")
	}
	if cfg.ImageName == "" {
		log.Fatal("IMAGE_NAME is required (env var or config.json)")
	}
	if cfg.GithubRepo == "" {
		log.Fatal("GITHUB_REPO is required (env var or config.json)")
	}

	cfg.ServerURL = strings.TrimRight(cfg.ServerURL, "/")
	return cfg
}

func handleRegister(w http.ResponseWriter, r *http.Request, cfg Config) {
	var req checkRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	body, _ := json.Marshal(map[string]string{
		"name":              req.Name,
		"email":             req.Email,
		"githubUsername":    req.GithubUsername,
		"dockerHubUsername": req.DockerHubUsername,
	})

	resp, err := http.Post(cfg.ServerURL+"/api/register", "application/json", bytes.NewReader(body))
	if err != nil {
		writeJSON(w, http.StatusBadGateway, map[string]string{"error": "failed to reach server: " + err.Error()})
		return
	}
	defer resp.Body.Close()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

func handleCheckGitHub(w http.ResponseWriter, r *http.Request, cfg Config) {
	var req checkRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	approved, err := ensureRegistered(cfg, req)
	if err != nil {
		writeJSON(w, http.StatusBadGateway, map[string]string{"error": err.Error()})
		return
	}
	if !approved {
		writeJSON(w, http.StatusOK, map[string]any{"status": "not_approved"})
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()

	checkErr := github.CheckRepoFile(ctx, req.GithubUsername, cfg.GithubRepo, "main.go")
	passed := checkErr == nil
	errMsg := ""
	if checkErr != nil {
		errMsg = checkErr.Error()
	}

	notifyErr := sendNotify(cfg, req, "github", passed, errMsg)
	if notifyErr != nil {
		writeJSON(w, http.StatusBadGateway, map[string]string{"error": "check ran but failed to notify server: " + notifyErr.Error()})
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"status":       boolToStatus(passed),
		"passed":       passed,
		"errorMessage": errMsg,
	})
}

func handleCheckDocker(w http.ResponseWriter, r *http.Request, cfg Config, runner *docker.Runner) {
	var req checkRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	approved, err := ensureRegistered(cfg, req)
	if err != nil {
		writeJSON(w, http.StatusBadGateway, map[string]string{"error": err.Error()})
		return
	}
	if !approved {
		writeJSON(w, http.StatusOK, map[string]any{"status": "not_approved"})
		return
	}

	imageRef := fmt.Sprintf("%s/%s", req.DockerHubUsername, cfg.ImageName)
	ctx, cancel := context.WithTimeout(r.Context(), 120*time.Second)
	defer cancel()

	result := runner.CheckStudent(ctx, imageRef, req.Email)

	notifyErr := sendNotify(cfg, req, "docker", result.Passed, result.ErrorMessage)
	if notifyErr != nil {
		writeJSON(w, http.StatusBadGateway, map[string]string{"error": "check ran but failed to notify server: " + notifyErr.Error()})
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"status":       boolToStatus(result.Passed),
		"passed":       result.Passed,
		"errorMessage": result.ErrorMessage,
	})
}

func ensureRegistered(cfg Config, req checkRequest) (bool, error) {
	body, _ := json.Marshal(map[string]string{
		"name":              req.Name,
		"email":             req.Email,
		"githubUsername":    req.GithubUsername,
		"dockerHubUsername": req.DockerHubUsername,
	})

	resp, err := http.Post(cfg.ServerURL+"/api/register", "application/json", bytes.NewReader(body))
	if err != nil {
		return false, fmt.Errorf("failed to reach server: %w", err)
	}
	defer resp.Body.Close()

	var result struct {
		Approved bool `json:"approved"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return false, fmt.Errorf("failed to parse register response: %w", err)
	}

	return result.Approved, nil
}

func sendNotify(cfg Config, req checkRequest, stage string, passed bool, errMsg string) error {
	body, _ := json.Marshal(map[string]any{
		"stage":             stage,
		"email":             req.Email,
		"name":              req.Name,
		"githubUsername":    req.GithubUsername,
		"dockerHubUsername": req.DockerHubUsername,
		"passed":            passed,
		"errorMessage":      errMsg,
	})

	resp, err := http.Post(cfg.ServerURL+"/api/notify", "application/json", bytes.NewReader(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("server returned %d: %s", resp.StatusCode, string(respBody))
	}
	return nil
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func boolToStatus(passed bool) string {
	if passed {
		return "passed"
	}
	return "failed"
}
