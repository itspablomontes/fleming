package middleware

import (
	"log/slog"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/itspablomontes/fleming/apps/backend/internal/config"
	"github.com/itspablomontes/fleming/apps/backend/internal/auth"
	"github.com/itspablomontes/fleming/apps/backend/internal/consent"
)

func AuthMiddleware(authService *auth.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		env := config.NormalizeEnv(os.Getenv("ENV"))
		overrideAddress := os.Getenv("DEV_OVERRIDE_WALLET_ADDRESS")
		if env == "dev" && overrideAddress != "" {
			slog.Debug("auth: using dev override", "address", overrideAddress)
			c.Set("user_address", overrideAddress)
			c.Next()
			return
		}

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

// ConsentMiddleware enforces patient-controlled access to medical data.
// It requires AuthMiddleware to have run first.
func ConsentMiddleware(consentService consent.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		userAddress, _ := c.Get("user_address")
		actor, ok := userAddress.(string)
		if !ok || actor == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}

		patientID := c.Query("patientId")
		if patientID == "" {
			patientID = actor
		}

		if actor == patientID {
			c.Set("target_patient", patientID)
			c.Next()
			return
		}

		permission := "read"
		if c.Request.Method == http.MethodPost || c.Request.Method == http.MethodPut || c.Request.Method == http.MethodDelete {
			permission = "write"
		}

		allowed, err := consentService.CheckPermission(c.Request.Context(), patientID, actor, permission)
		if err != nil {
			slog.Error("consent check error", "actor", actor, "patient", patientID, "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to verify access permissions"})
			c.Abort()
			return
		}

		if !allowed {
			slog.Warn("access denied: no valid consent", "actor", actor, "patient", patientID)
			c.JSON(http.StatusForbidden, gin.H{"error": "access denied: you do not have permission to access this patient's data"})
			c.Abort()
			return
		}

		c.Set("target_patient", patientID)
		c.Next()
	}
}
