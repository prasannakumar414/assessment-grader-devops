package handlers

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"docker-workshop-assesment-grader/internal/models"
)

type StudentHandler struct {
	DB *gorm.DB
}

type studentPayload struct {
	Name              string `json:"name" binding:"required"`
	Email             string `json:"email" binding:"required,email"`
	GitHubUsername    string `json:"githubUsername" binding:"required"`
	DockerHubUsername string `json:"dockerHubUsername" binding:"required"`
}

func (h *StudentHandler) CreateStudent(c *gin.Context) {
	var payload studentPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	student := models.Student{
		Name:              strings.TrimSpace(payload.Name),
		Email:             strings.ToLower(strings.TrimSpace(payload.Email)),
		GitHubUsername:    strings.TrimSpace(payload.GitHubUsername),
		DockerHubUsername: strings.TrimSpace(payload.DockerHubUsername),
		Approved:          true,
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

	c.JSON(http.StatusCreated, student)
}

func (h *StudentHandler) ListStudents(c *gin.Context) {
	var students []models.Student
	query := h.DB.Order("created_at desc")

	if status := c.Query("status"); status != "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "filtering by single status is no longer supported; use approved filter"})
		return
	}
	if approved := c.Query("approved"); approved != "" {
		query = query.Where("approved = ?", approved == "true")
	}

	if err := query.Find(&students).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list students"})
		return
	}

	c.JSON(http.StatusOK, students)
}

func (h *StudentHandler) GetStudent(c *gin.Context) {
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

	c.JSON(http.StatusOK, student)
}

func (h *StudentHandler) UpdateStudent(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}

	var payload studentPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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

	student.Name = strings.TrimSpace(payload.Name)
	student.Email = strings.ToLower(strings.TrimSpace(payload.Email))
	student.GitHubUsername = strings.TrimSpace(payload.GitHubUsername)
	student.DockerHubUsername = strings.TrimSpace(payload.DockerHubUsername)

	if err := h.DB.Save(&student).Error; err != nil {
		if isUniqueConstraintError(err) {
			c.JSON(http.StatusConflict, gin.H{"error": "email already exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update student"})
		return
	}

	c.JSON(http.StatusOK, student)
}

func (h *StudentHandler) DeleteStudent(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}

	result := h.DB.Delete(&models.Student{}, id)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete student"})
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "student not found"})
		return
	}

	c.Status(http.StatusNoContent)
}

func parseID(c *gin.Context) (uint, bool) {
	rawID := c.Param("id")
	parsed, err := strconv.ParseUint(rawID, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid student id"})
		return 0, false
	}
	return uint(parsed), true
}

func isUniqueConstraintError(err error) bool {
	return strings.Contains(strings.ToLower(err.Error()), "unique")
}
