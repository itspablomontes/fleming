package api

import (
	"context"
	"log/slog"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/itspablomontes/fleming/api/internal/auth"
	"github.com/itspablomontes/fleming/api/internal/middleware"
	"github.com/itspablomontes/fleming/api/internal/timeline"
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
	authService := auth.NewService(authRepo, jwtSecret)

	timelineRepo := timeline.NewRepository(db)
	timelineService := timeline.NewService(timelineRepo)

	authService.StartCleanup(context.Background())

	authHandler := auth.NewHandler(authService)
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
	timelineHandler.RegisterRoutes(apiGroup)

	return r
}
