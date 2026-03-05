package github

import (
	"context"
	"fmt"
	"net/http"
)

func CheckRepoFile(ctx context.Context, username, repo, filePath string) error {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/contents/%s", username, repo, filePath)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("github API request failed: %w", err)
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		return nil
	case http.StatusNotFound:
		return fmt.Errorf("file %q not found in repo %s/%s", filePath, username, repo)
	default:
		return fmt.Errorf("github API returned status %d for %s/%s/%s", resp.StatusCode, username, repo, filePath)
	}
}
