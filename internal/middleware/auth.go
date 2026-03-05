package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"docker-workshop-assesment-grader/internal/auth"
)

func RequireAuth(sessions *auth.SessionStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr := extractToken(c)
		if tokenStr == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing authentication token"})
			return
		}

		username, ok := sessions.Validate(tokenStr)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired session"})
			return
		}

		c.Set("admin_user", username)
		c.Next()
	}
}

func extractToken(c *gin.Context) string {
	if header := c.GetHeader("Authorization"); header != "" {
		parts := strings.SplitN(header, " ", 2)
		if len(parts) == 2 && strings.EqualFold(parts[0], "bearer") {
			return strings.TrimSpace(parts[1])
		}
	}

	// Fallback: ?token=<token> query param (used by EventSource which can't set headers)
	if t := c.Query("token"); t != "" {
		return t
	}

	return ""
}
