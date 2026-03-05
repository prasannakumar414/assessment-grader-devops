package handlers

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"docker-workshop-assesment-grader/internal/models"
	"docker-workshop-assesment-grader/internal/sse"
)

type NotifyHandler struct {
	DB  *gorm.DB
	Hub *sse.Hub
}

type notifyPayload struct {
	Stage             string `json:"stage" binding:"required,oneof=github docker k8s"`
	Email             string `json:"email" binding:"required,email"`
	Name              string `json:"name"`
	GitHubUsername    string `json:"githubUsername"`
	DockerHubUsername string `json:"dockerHubUsername"`
	Passed            bool   `json:"passed"`
	ErrorMessage      string `json:"errorMessage"`
}

func (h *NotifyHandler) Notify(c *gin.Context) {
	var payload notifyPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	email := strings.ToLower(strings.TrimSpace(payload.Email))

	var student models.Student
	if err := h.DB.Where("email = ?", email).First(&student).Error; err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "student not registered"})
		return
	}

	if !student.Approved {
		c.JSON(http.StatusForbidden, gin.H{"error": "student not approved"})
		return
	}

	now := time.Now()
	status := models.StatusFailed
	if payload.Passed {
		status = models.StatusPassed
	}

	wasAlreadyPassed := stageStatus(&student, payload.Stage) == models.StatusPassed

	setStageResult(&student, payload.Stage, status, payload.ErrorMessage, &now)

	if err := h.DB.Save(&student).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save status"})
		return
	}

	if payload.Passed && !wasAlreadyPassed {
		h.Hub.Broadcast(sse.Event{
			Type: "stage_complete",
			Data: map[string]any{
				"studentName": student.Name,
				"stageName":   payload.Stage,
			},
		})

		if student.AllPassed() {
			h.Hub.Broadcast(sse.Event{
				Type: "all_complete",
				Data: map[string]any{
					"studentName": student.Name,
				},
			})
		}
	}

	c.JSON(http.StatusOK, student)
}

func stageStatus(s *models.Student, stage string) string {
	switch stage {
	case models.StageGitHub:
		return s.GitHubStatus
	case models.StageDocker:
		return s.DockerStatus
	case models.StageK8s:
		return s.K8sStatus
	}
	return ""
}

func setStageResult(s *models.Student, stage, status, errMsg string, checkedAt *time.Time) {
	if status == models.StatusPassed {
		errMsg = ""
	}
	switch stage {
	case models.StageGitHub:
		s.GitHubStatus = status
		s.GitHubErrorMessage = errMsg
		s.GitHubLastCheckedAt = checkedAt
	case models.StageDocker:
		s.DockerStatus = status
		s.DockerErrorMessage = errMsg
		s.DockerLastCheckedAt = checkedAt
	case models.StageK8s:
		s.K8sStatus = status
		s.K8sErrorMessage = errMsg
		s.K8sLastCheckedAt = checkedAt
	}
}
