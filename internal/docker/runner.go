package docker

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"runtime"
	"strings"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

type Runner struct {
	Client       *client.Client
	HTTPClient   *http.Client
	ReadyTimeout time.Duration
	VerifyHost   string
}

type CheckResult struct {
	Passed       bool
	ErrorMessage string
}

func NewRunner(verifyHost string) (*Runner, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, fmt.Errorf("create docker client: %w", err)
	}

	if verifyHost == "" {
		verifyHost = "127.0.0.1"
	}

	return &Runner{
		Client:       cli,
		HTTPClient:   &http.Client{Timeout: 5 * time.Second},
		ReadyTimeout: 30 * time.Second,
		VerifyHost:   verifyHost,
	}, nil
}

func (r *Runner) CheckStudent(ctx context.Context, imageName string, expectedEmail string) CheckResult {
	imageRef := normalizeImage(imageName)

	if err := r.pullImage(ctx, imageRef); err != nil {
		return CheckResult{Passed: false, ErrorMessage: err.Error()}
	}

	containerID, verifyURL, err := r.startContainer(ctx, imageRef)
	if err != nil {
		return CheckResult{Passed: false, ErrorMessage: err.Error()}
	}

	defer r.cleanupContainer(context.Background(), containerID)

	log.Printf("verifying student container at %s", verifyURL)

	passed, err := r.verifyContainer(ctx, verifyURL, expectedEmail)
	if err != nil {
		return CheckResult{Passed: false, ErrorMessage: err.Error()}
	}

	if !passed {
		return CheckResult{
			Passed:       false,
			ErrorMessage: "container response email does not match registered student email",
		}
	}

	return CheckResult{Passed: true}
}

func (r *Runner) pullImage(ctx context.Context, imageRef string) error {
	options := image.PullOptions{
		Platform: "linux/" + runtime.GOARCH,
	}
	reader, err := r.Client.ImagePull(ctx, imageRef, options)
	if err != nil {
		return fmt.Errorf("docker pull failed for %s: %w", imageRef, err)
	}
	defer reader.Close()

	_, _ = io.Copy(io.Discard, reader)
	return nil
}

// startContainer creates and starts the student container, returning its ID
// and the URL to use for verification.
// When VerifyHost is "container", it uses the container's bridge network IP
// directly (reliable for container-to-container communication).
// Otherwise, it uses host port mapping with the configured VerifyHost.
func (r *Runner) startContainer(ctx context.Context, imageRef string) (string, string, error) {
	containerPort, err := nat.NewPort("tcp", "8080")
	if err != nil {
		return "", "", fmt.Errorf("create container port mapping: %w", err)
	}

	config := &container.Config{
		Image: imageRef,
		ExposedPorts: nat.PortSet{
			containerPort: struct{}{},
		},
	}

	hostConfig := &container.HostConfig{
		PortBindings: nat.PortMap{
			containerPort: []nat.PortBinding{
				{
					HostIP:   "0.0.0.0",
					HostPort: "0",
				},
			},
		},
	}

	resp, err := r.Client.ContainerCreate(ctx, config, hostConfig, &network.NetworkingConfig{}, nil, "")
	if err != nil {
		return "", "", fmt.Errorf("create container: %w", err)
	}

	if err := r.Client.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		return "", "", fmt.Errorf("start container: %w", err)
	}

	if r.VerifyHost == "container" {
		ip, ipErr := r.resolveContainerIP(ctx, resp.ID)
		if ipErr != nil {
			return "", "", ipErr
		}
		verifyURL := fmt.Sprintf("http://%s:8080/api/info", ip)
		return resp.ID, verifyURL, nil
	}

	hostPort, portErr := r.resolveHostPort(ctx, resp.ID, containerPort)
	if portErr != nil {
		return "", "", portErr
	}
	verifyURL := fmt.Sprintf("http://%s:%s/api/info", r.VerifyHost, hostPort)
	return resp.ID, verifyURL, nil
}

func (r *Runner) resolveContainerIP(ctx context.Context, containerID string) (string, error) {
	deadline := time.Now().Add(10 * time.Second)
	for time.Now().Before(deadline) {
		inspection, err := r.Client.ContainerInspect(ctx, containerID)
		if err != nil {
			return "", fmt.Errorf("inspect container: %w", err)
		}

		if ip := inspection.NetworkSettings.IPAddress; ip != "" {
			return ip, nil
		}

		// Also check named networks
		for _, net := range inspection.NetworkSettings.Networks {
			if net.IPAddress != "" {
				return net.IPAddress, nil
			}
		}

		time.Sleep(250 * time.Millisecond)
	}

	return "", errors.New("timed out resolving container IP address")
}

func (r *Runner) resolveHostPort(ctx context.Context, containerID string, containerPort nat.Port) (string, error) {
	deadline := time.Now().Add(10 * time.Second)
	for time.Now().Before(deadline) {
		inspection, err := r.Client.ContainerInspect(ctx, containerID)
		if err != nil {
			return "", fmt.Errorf("inspect container: %w", err)
		}

		bindings, ok := inspection.NetworkSettings.Ports[containerPort]
		if ok && len(bindings) > 0 && bindings[0].HostPort != "" {
			return bindings[0].HostPort, nil
		}

		time.Sleep(250 * time.Millisecond)
	}

	return "", errors.New("timed out resolving mapped host port")
}

func (r *Runner) verifyContainer(ctx context.Context, url string, expectedEmail string) (bool, error) {
	deadline := time.Now().Add(r.ReadyTimeout)

	for time.Now().Before(deadline) {
		select {
		case <-ctx.Done():
			return false, ctx.Err()
		default:
		}

		resp, err := r.HTTPClient.Get(url)
		if err != nil {
			time.Sleep(1 * time.Second)
			continue
		}

		body, readErr := io.ReadAll(resp.Body)
		resp.Body.Close()
		if readErr != nil {
			time.Sleep(1 * time.Second)
			continue
		}

		if resp.StatusCode >= http.StatusBadRequest {
			time.Sleep(1 * time.Second)
			continue
		}

		if matchesExpectedEmail(body, expectedEmail) {
			return true, nil
		}
		return false, nil
	}

	return false, fmt.Errorf("timed out waiting for response from %s", url)
}

func (r *Runner) cleanupContainer(ctx context.Context, containerID string) {
	timeout := 3
	_ = r.Client.ContainerStop(ctx, containerID, container.StopOptions{Timeout: &timeout})
	_ = r.Client.ContainerRemove(ctx, containerID, container.RemoveOptions{Force: true})
}

func normalizeImage(rollNo string) string {
	image := strings.TrimSpace(rollNo)
	if !strings.Contains(image, ":") {
		return image + ":latest"
	}
	return image
}

func matchesExpectedEmail(body []byte, expectedEmail string) bool {
	expected := strings.ToLower(strings.TrimSpace(expectedEmail))

	var parsed map[string]any
	if err := json.Unmarshal(body, &parsed); err == nil {
		if emailRaw, ok := parsed["email"]; ok {
			return strings.EqualFold(fmt.Sprint(emailRaw), expected)
		}
	}

	plain := strings.ToLower(strings.TrimSpace(string(body)))
	return plain == expected
}
