package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"docker-workshop-assesment-grader/internal/models"
)

type ApprovalHandler struct {
	DB *gorm.DB
}

func (h *ApprovalHandler) ApproveOne(c *gin.Context) {
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

	student.Approved = true
	if err := h.DB.Save(&student).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to approve student"})
		return
	}

	c.JSON(http.StatusOK, student)
}

func (h *ApprovalHandler) ApproveAll(c *gin.Context) {
	result := h.DB.Model(&models.Student{}).Where("approved = ?", false).Update("approved", true)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to approve students"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"approved": result.RowsAffected,
	})
}
