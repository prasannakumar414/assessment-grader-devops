/*
Copyright 2026.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"

	examplecomv1 "github.com/Kvnpsiddhartha/health-operator/api/v1"
)

// AppHealthReconciler reconciles an AppHealth object (validates student deployments).
type AppHealthReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=example.com,resources=apphealths,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=example.com,resources=apphealths/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=example.com,resources=apphealths/finalizers,verbs=update
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch
// +kubebuilder:rbac:groups="",resources=services;pods,verbs=get;list;watch

// Reconcile validates that the student has correctly deployed their app (Deployment + Service)
// and that the app serves /info with an email for grading.
func (r *AppHealthReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := logf.FromContext(ctx)

	appHealth := &examplecomv1.AppHealth{}
	if err := r.Get(ctx, req.NamespacedName, appHealth); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	deploymentStatus := "Not Found"
	serviceStatus := "Not Found"
	imageValid := true
	email := "Error"

	// Fetch Deployment (apps/v1)
	deployment := &appsv1.Deployment{}
	err := r.Get(ctx, client.ObjectKey{Namespace: appHealth.Spec.Namespace, Name: appHealth.Spec.AppName}, deployment)
	if err != nil {
		deploymentStatus = "Unhealthy"
	} else if deployment.Status.ReadyReplicas > 0 {
		deploymentStatus = "Healthy"
		if appHealth.Spec.ExpectedImage != "" {
			if len(deployment.Spec.Template.Spec.Containers) > 0 {
				imageValid = deployment.Spec.Template.Spec.Containers[0].Image == appHealth.Spec.ExpectedImage
			} else {
				imageValid = false
			}
		}
	} else {
		deploymentStatus = "Unhealthy"
	}

	// Fetch Service
	service := &corev1.Service{}
	err = r.Get(ctx, client.ObjectKey{Namespace: appHealth.Spec.Namespace, Name: appHealth.Spec.AppName}, service)
	if err != nil {
		serviceStatus = "Unhealthy"
	} else if len(service.Spec.Ports) > 0 {
		serviceStatus = "Healthy"
	} else {
		serviceStatus = "Unhealthy"
	}

	// Call the service API to get the email only when service is reachable
	infoPath := appHealth.Spec.InfoPath
	if infoPath == "" {
		infoPath = "/info"
	}
	if serviceStatus == "Healthy" {
		email, err = getEmailFromService(ctx, appHealth.Spec.Namespace, appHealth.Spec.AppName, infoPath)
		if err != nil {
			logger.Error(err, "failed to retrieve email from service")
			email = "Error"
		}
	}

	healthy := deploymentStatus == "Healthy" && serviceStatus == "Healthy" && imageValid

	// Update status
	appHealth.Status.DeploymentStatus = deploymentStatus
	appHealth.Status.ServiceStatus = serviceStatus
	appHealth.Status.ImageValid = imageValid
	appHealth.Status.Healthy = healthy
	appHealth.Status.Email = email
	appHealth.Status.LastChecked = time.Now().Format(time.RFC3339)
	if healthy {
		appHealth.Status.Message = "Deployment validated successfully"
	} else {
		appHealth.Status.Message = fmt.Sprintf("Deployment=%s, Service=%s, ImageValid=%v", deploymentStatus, serviceStatus, imageValid)
	}

	logger.Info("Updating AppHealth status", "status", appHealth.Status)

	if err := r.Status().Update(ctx, appHealth); err != nil {
		logger.Error(err, "unable to update AppHealth status")
		return ctrl.Result{}, err
	}

	// Report to external API if configured
	if appHealth.Spec.ReportURL != "" {
		reportToAPIServer(appHealth)
	}

	return ctrl.Result{
		RequeueAfter: 5 * time.Minute,
	}, nil
}

// getEmailFromService calls the deployed service's info endpoint.
func getEmailFromService(ctx context.Context, namespace, appName, infoPath string) (string, error) {
	url := fmt.Sprintf("http://%s.%s.svc.cluster.local%s", appName, namespace, infoPath)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to call service API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("service API returned status %d", resp.StatusCode)
	}

	var responseBody struct {
		Email string `json:"email"`
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}
	if err := json.Unmarshal(body, &responseBody); err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %w", err)
	}
	return responseBody.Email, nil
}

// reportToAPIServer sends validation result to the grading server's /api/notify endpoint.
func reportToAPIServer(appHealth *examplecomv1.AppHealth) {
	var errorMsg string
	if !appHealth.Status.Healthy {
		errorMsg = appHealth.Status.Message
	}

	body := map[string]interface{}{
		"email":        appHealth.Status.Email,
		"passed":       appHealth.Status.Healthy,
		"stage":        "k8s",
		"errorMessage": errorMsg,
	}
	jsonData, err := json.Marshal(body)
	if err != nil {
		return
	}

	url := strings.TrimRight(appHealth.Spec.ReportURL, "/") + "/api/notify"
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return
	}
	defer resp.Body.Close()
}

// SetupWithManager sets up the controller with the Manager.
func (r *AppHealthReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&examplecomv1.AppHealth{}).
		Named("apphealth").
		Complete(r)
}
