package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"docker-workshop-assesment-grader/internal/models"
	"docker-workshop-assesment-grader/internal/sse"
)

type RegisterHandler struct {
	DB  *gorm.DB
	Hub *sse.Hub
}

type registerPayload struct {
	Name              string `json:"name" binding:"required"`
	Email             string `json:"email" binding:"required,email"`
	GitHubUsername    string `json:"githubUsername" binding:"required"`
	DockerHubUsername string `json:"dockerHubUsername" binding:"required"`
}

func (h *RegisterHandler) Register(c *gin.Context) {
	var payload registerPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	email := strings.ToLower(strings.TrimSpace(payload.Email))

	var student models.Student
	result := h.DB.Where("email = ?", email).First(&student)

	if result.Error == nil {
		c.JSON(http.StatusOK, student)
		return
	}

	student = models.Student{
		Name:              strings.TrimSpace(payload.Name),
		Email:             email,
		GitHubUsername:    strings.TrimSpace(payload.GitHubUsername),
		DockerHubUsername: strings.TrimSpace(payload.DockerHubUsername),
		Approved:          false,
		GitHubStatus:      models.StatusPending,
		DockerStatus:      models.StatusPending,
		K8sStatus:         models.StatusPending,
	}

	if err := h.DB.Create(&student).Error; err != nil {
		if isUniqueConstraintError(err) {
			c.JSON(http.StatusConflict, gin.H{"error": "email already exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create student"})
		return
	}

	h.Hub.Broadcast(sse.Event{
		Type: "new_registration",
		Data: map[string]any{
			"studentName": student.Name,
			"studentId":   student.ID,
		},
	})

	c.JSON(http.StatusCreated, student)
}
