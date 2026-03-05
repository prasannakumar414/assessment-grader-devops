package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"docker-workshop-assesment-grader/internal/auth"
)

type AuthHandler struct {
	Username string
	Password string
	Sessions *auth.SessionStore
}

type loginPayload struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *AuthHandler) Login(c *gin.Context) {
	var payload loginPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username and password are required"})
		return
	}

	if payload.Username != h.Username || payload.Password != h.Password {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	token, err := h.Sessions.Create(payload.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create session"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
