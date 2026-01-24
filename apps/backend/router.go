package api

import (
	"context"
	"log/slog"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/itspablomontes/fleming/apps/backend/internal/audit"
	"github.com/itspablomontes/fleming/apps/backend/internal/auth"
	"github.com/itspablomontes/fleming/apps/backend/internal/consent"
	"github.com/itspablomontes/fleming/apps/backend/internal/middleware"
	"github.com/itspablomontes/fleming/apps/backend/internal/timeline"
	"gorm.io/gorm"
)

func NewRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()

	jwtSecret := os.Getenv("JWT_SECRET")
	env := os.Getenv("ENV")
	if jwtSecret == "" {
		if env == "production" || env == "staging" {
			slog.Error("JWT_SECRET is required in production/staging environments")
			os.Exit(1)
		}
		jwtSecret = "dev-secret-do-not-use-in-prod"
		slog.Warn("JWT_SECRET not set, using insecure default for development", "env", env)
	} else {
		slog.Info("JWT_SECRET loaded", "length", len(jwtSecret), "env", env)
	}

	authRepo := auth.NewGormRepository(db)
	auditRepo := audit.NewRepository(db)
	consentRepo := consent.NewRepository(db)
	timelineRepo := timeline.NewRepository(db)

	auditService := audit.NewService(auditRepo)
	consentService := consent.NewService(consentRepo)
	authService := auth.NewService(authRepo, jwtSecret, auditService)
	timelineService := timeline.NewService(timelineRepo, auditService)

	authService.StartCleanup(context.Background())

	authHandler := auth.NewHandler(authService)
	auditHandler := audit.NewHandler(auditService)
	consentHandler := consent.NewHandler(consentService)
	timelineHandler := timeline.NewHandler(timelineService)

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"service": "fleming-backend",
		})
	})

	authHandler.RegisterRoutes(r.Group("/auth"))

	r.GET("/auth/me", middleware.AuthMiddleware(authService), authHandler.HandleMe)

	apiGroup := r.Group("/api")
	apiGroup.Use(middleware.AuthMiddleware(authService))

	auditHandler.RegisterRoutes(apiGroup)
	consentHandler.RegisterRoutes(apiGroup)

	// Timeline routes are protected by both Auth and Consent middleware
	timelineGroup := apiGroup.Group("")
	timelineGroup.Use(middleware.ConsentMiddleware(consentService))
	timelineHandler.RegisterRoutes(timelineGroup)

	return r
}
