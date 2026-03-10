package handlers

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"docker-workshop-assesment-grader/internal/docker"
	"docker-workshop-assesment-grader/internal/models"
	"docker-workshop-assesment-grader/internal/sse"
)

type DockerCheckHandler struct {
	DB        *gorm.DB
	Hub       *sse.Hub
	Runner    *docker.Runner
	ImageName string
}

func (h *DockerCheckHandler) Check(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid student id"})
		return
	}

	var student models.Student
	if err := h.DB.First(&student, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "student not found"})
		return
	}

	if !student.Approved {
		c.JSON(http.StatusForbidden, gin.H{"error": "student not approved"})
		return
	}

	imageRef := fmt.Sprintf("%s/%s", student.DockerHubUsername, h.ImageName)
	ctx, cancel := context.WithTimeout(c.Request.Context(), 120*time.Second)
	defer cancel()

	result := h.Runner.CheckStudent(ctx, imageRef, student.Email)

	now := time.Now()
	status := models.StatusFailed
	if result.Passed {
		status = models.StatusPassed
	}

	wasAlreadyPassed := student.DockerStatus == models.StatusPassed
	setStageResult(&student, models.StageDocker, status, result.ErrorMessage, &now)

	if err := h.DB.Save(&student).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save status"})
		return
	}

	if result.Passed && !wasAlreadyPassed {
		h.Hub.Broadcast(sse.Event{
			Type: "stage_complete",
			Data: map[string]any{
				"studentName": student.Name,
				"stageName":   models.StageDocker,
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

	c.JSON(http.StatusOK, gin.H{
		"passed":       result.Passed,
		"errorMessage": result.ErrorMessage,
	})
}
