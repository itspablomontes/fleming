package middleware

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/itspablomontes/fleming/api/internal/auth"
)

func AuthMiddleware(authService *auth.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		source := "header"

		if authHeader == "" {
			cookie, err := c.Cookie("auth_token")
			if err == nil {
				authHeader = "Bearer " + cookie
				source = "cookie"
			}
		}

		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			slog.Debug("auth: no token found", "hasHeader", authHeader != "")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		address, err := authService.ValidateJWT(tokenString)
		if err != nil {
			slog.Warn("auth: token validation failed", "source", source, "error", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			c.Abort()
			return
		}

		slog.Debug("auth: success", "address", address, "source", source)
		c.Set("user_address", address)
		c.Next()
	}
}
