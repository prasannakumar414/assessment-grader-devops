package handlers

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"docker-workshop-assesment-grader/internal/docker"
	"docker-workshop-assesment-grader/internal/models"
)

type CheckHandler struct {
	DB     *gorm.DB
	Runner *docker.Runner
}

type checkResponse struct {
	StudentID    uint   `json:"studentId"`
	RollNo       string `json:"rollNo"`
	Status       string `json:"status"`
	ErrorMessage string `json:"errorMessage,omitempty"`
}

func (h *CheckHandler) RunCheckAll(c *gin.Context) {
	var students []models.Student
	if err := h.DB.Where("status <> ?", models.StatusPassed).Find(&students).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list students for check"})
		return
	}

	results := make([]checkResponse, 0, len(students))
	for i := range students {
		student := students[i]
		result := h.runSingle(c, &student)
		results = append(results, result)
	}

	c.JSON(http.StatusOK, gin.H{
		"count":   len(results),
		"results": results,
	})
}

func (h *CheckHandler) RunCheckByID(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}

	var student models.Student
	if err := h.DB.First(&student, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "student not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get student"})
		return
	}

	result := h.runSingle(c, &student)
	c.JSON(http.StatusOK, result)
}

func (h *CheckHandler) runSingle(c *gin.Context, student *models.Student) checkResponse {
	now := time.Now()
	check := h.Runner.CheckStudent(c.Request.Context(), student.RollNo, student.Email)

	student.LastCheckedAt = &now
	student.ErrorMessage = check.ErrorMessage
	if check.Passed {
		student.Status = models.StatusPassed
		student.ErrorMessage = ""
	} else {
		student.Status = models.StatusFailed
	}

	if err := h.DB.Save(student).Error; err != nil {
		return checkResponse{
			StudentID:    student.ID,
			RollNo:       student.RollNo,
			Status:       models.StatusFailed,
			ErrorMessage: "failed to persist check status: " + err.Error(),
		}
	}

	return checkResponse{
		StudentID:    student.ID,
		RollNo:       student.RollNo,
		Status:       student.Status,
		ErrorMessage: student.ErrorMessage,
	}
}
