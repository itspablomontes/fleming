package api

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/itspablomontes/fleming/apps/backend/internal/audit"
	"github.com/itspablomontes/fleming/apps/backend/internal/auth"
	"github.com/itspablomontes/fleming/apps/backend/internal/config"
	"github.com/itspablomontes/fleming/apps/backend/internal/consent"
	"github.com/itspablomontes/fleming/apps/backend/internal/middleware"
	"github.com/itspablomontes/fleming/apps/backend/internal/storage"
	"github.com/itspablomontes/fleming/apps/backend/internal/timeline"
	"gorm.io/gorm"
)

func NewRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()

	// Enable CORS
	r.Use(func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		if origin != "" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		}

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	})

	jwtSecret := os.Getenv("JWT_SECRET")
	env := config.NormalizeEnv(os.Getenv("ENV"))
	if jwtSecret == "" {
		if config.IsProductionLike(env) {
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

	storageEndpointRaw := firstNonEmpty(os.Getenv("STORAGE_ENDPOINT"), os.Getenv("S3_ENDPOINT"))
	storageAccessKey := firstNonEmpty(os.Getenv("STORAGE_ACCESS_KEY"), os.Getenv("S3_ACCESS_KEY"))
	storageSecretKey := firstNonEmpty(os.Getenv("STORAGE_SECRET_KEY"), os.Getenv("S3_SECRET_KEY"))
	storageBucket := firstNonEmpty(os.Getenv("STORAGE_BUCKET"), os.Getenv("S3_BUCKET"))

	storageUseSSLStr := firstNonEmpty(os.Getenv("STORAGE_USE_SSL"), os.Getenv("S3_SSL"))
	storageUseSSL, hasStorageUseSSL, err := parseOptionalBool(storageUseSSLStr)
	if err != nil {
		slog.Error("Invalid STORAGE_USE_SSL/S3_SSL value", "value", storageUseSSLStr, "error", err)
		os.Exit(1)
	}

	if storageEndpointRaw == "" {
		if config.IsProductionLike(env) {
			slog.Error("STORAGE_ENDPOINT (or S3_ENDPOINT) is required in production/staging")
			os.Exit(1)
		}
		slog.Warn("STORAGE_ENDPOINT not set; defaulting to localhost:9000 for development")
		storageEndpointRaw = "localhost:9000"
	}
	if storageAccessKey == "" {
		if config.IsProductionLike(env) {
			slog.Error("STORAGE_ACCESS_KEY (or S3_ACCESS_KEY) is required in production/staging")
			os.Exit(1)
		}
		storageAccessKey = "minioadmin"
	}
	if storageSecretKey == "" {
		if config.IsProductionLike(env) {
			slog.Error("STORAGE_SECRET_KEY (or S3_SECRET_KEY) is required in production/staging")
			os.Exit(1)
		}
		storageSecretKey = "minioadmin"
	}
	if storageBucket == "" {
		if config.IsProductionLike(env) {
			slog.Error("STORAGE_BUCKET (or S3_BUCKET) is required in production/staging")
			os.Exit(1)
		}
		storageBucket = "fleming-blobs"
	}

	storageEndpoint, inferredSSL, err := normalizeStorageEndpoint(storageEndpointRaw)
	if err != nil {
		slog.Error("Invalid STORAGE_ENDPOINT/S3_ENDPOINT", "value", storageEndpointRaw, "error", err)
		os.Exit(1)
	}
	if !hasStorageUseSSL {
		if inferredSSL != nil {
			storageUseSSL = *inferredSSL
		} else if config.IsProductionLike(env) {
			slog.Error("STORAGE_USE_SSL (or S3_SSL) is required in production/staging when endpoint has no scheme")
			os.Exit(1)
		}
	}

	storageService, err := storage.NewMinIOStorage(storageEndpoint, storageAccessKey, storageSecretKey, storageUseSSL)
	if err != nil {
		slog.Error("Failed to initialize storage service", "error", err)
		os.Exit(1)
	}

	auditService := audit.NewService(auditRepo)
	consentService := consent.NewService(consentRepo, auditService)
	authService := auth.NewService(authRepo, jwtSecret, auditService)
	timelineService := timeline.NewService(timelineRepo, auditService, storageService, storageBucket)

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

	authGroup := r.Group("/api/auth")
	authHandler.RegisterRoutes(authGroup)

	r.GET("/api/auth/me", middleware.AuthMiddleware(authService), authHandler.HandleMe)

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

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if strings.TrimSpace(v) != "" {
			return strings.TrimSpace(v)
		}
	}
	return ""
}

func parseOptionalBool(v string) (value bool, ok bool, err error) {
	if strings.TrimSpace(v) == "" {
		return false, false, nil
	}
	switch strings.ToLower(strings.TrimSpace(v)) {
	case "true", "1", "yes", "y":
		return true, true, nil
	case "false", "0", "no", "n":
		return false, true, nil
	default:
		return false, false, fmt.Errorf("invalid boolean %q", v)
	}
}

// normalizeStorageEndpoint accepts either a host[:port] (recommended) or an http(s) URL.
// It returns the host[:port] suitable for minio-go, plus an inferred TLS value when a scheme was provided.
func normalizeStorageEndpoint(raw string) (hostPort string, inferredSSL *bool, err error) {
	r := strings.TrimSpace(raw)
	if r == "" {
		return "", nil, fmt.Errorf("empty endpoint")
	}

	if strings.HasPrefix(r, "http://") || strings.HasPrefix(r, "https://") {
		u, err := url.Parse(r)
		if err != nil {
			return "", nil, fmt.Errorf("parse url: %w", err)
		}
		if u.Host == "" {
			return "", nil, fmt.Errorf("missing host in %q", r)
		}
		ssl := u.Scheme == "https"
		return u.Host, &ssl, nil
	}

	return r, nil, nil
}
