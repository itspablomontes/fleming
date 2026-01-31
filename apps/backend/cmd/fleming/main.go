package main

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	api "github.com/itspablomontes/fleming/apps/backend"
	"github.com/itspablomontes/fleming/apps/backend/internal/config"
	"github.com/itspablomontes/fleming/apps/backend/internal/audit"
	"github.com/itspablomontes/fleming/apps/backend/internal/auth"
	"github.com/itspablomontes/fleming/apps/backend/internal/consent"
	"github.com/itspablomontes/fleming/apps/backend/internal/timeline"
)

func main() {
	env := config.NormalizeEnv(os.Getenv("ENV"))
	logLevel := slog.LevelDebug
	if config.IsProductionLike(env) {
		logLevel = slog.LevelInfo
	}
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: logLevel,
	})))

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		slog.Error("DATABASE_URL not set")
		os.Exit(1)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		slog.Warn("PORT not set; defaulting to 8080")
	}

	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})
	if err != nil {
		slog.Error("Failed to connect to database", "error", err)
		os.Exit(1)
	}

	sqlDB, err := db.DB()
	if err != nil {
		slog.Error("Failed to get generic database object", "error", err)
		os.Exit(1)
	}
	defer sqlDB.Close()

	if err := sqlDB.Ping(); err != nil {
		slog.Error("failed to ping database", "error", err)
		os.Exit(1)
	}
	slog.Info("database connected")

	// Connection pool tuning (important for serverless Postgres like Neon).
	applyConnPoolSettings(sqlDB, env)

	if err := db.AutoMigrate(
		&auth.Challenge{},
		&auth.User{},
		&timeline.TimelineEvent{},
		&timeline.EventEdge{},
		&timeline.EventFile{},
		&timeline.EventFileAccess{},
		&audit.AuditEntry{},
		&audit.AuditBatch{},
		&consent.ConsentGrant{},
	); err != nil {
		slog.Error("failed to auto-migrate schema", "error", err)
		os.Exit(1)
	}

	router := api.NewRouter(db)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	go func() {
		slog.Info("Starting server", "port", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("Failed to start server", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	slog.Info("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("server forced to shutdown", "error", err)
		os.Exit(1)
	}

	slog.Info("server exiting")
}

func applyConnPoolSettings(sqlDB *sql.DB, env string) {
	// Defaults only for prod-like environments to reduce accidental connection storms.
	defaultMaxOpen := 5
	defaultMaxIdle := 2
	defaultMaxLifetime := 30 * time.Minute

	maxOpen, maxOpenSet, err := getOptionalIntEnv("DB_MAX_OPEN_CONNS")
	if err != nil {
		slog.Error("invalid DB_MAX_OPEN_CONNS", "error", err)
		os.Exit(1)
	}
	maxIdle, maxIdleSet, err := getOptionalIntEnv("DB_MAX_IDLE_CONNS")
	if err != nil {
		slog.Error("invalid DB_MAX_IDLE_CONNS", "error", err)
		os.Exit(1)
	}
	maxLifetime, maxLifetimeSet, err := getOptionalDurationEnv("DB_CONN_MAX_LIFETIME")
	if err != nil {
		slog.Error("invalid DB_CONN_MAX_LIFETIME", "error", err)
		os.Exit(1)
	}

	if config.IsProductionLike(env) {
		if !maxOpenSet {
			maxOpen, maxOpenSet = defaultMaxOpen, true
		}
		if !maxIdleSet {
			maxIdle, maxIdleSet = defaultMaxIdle, true
		}
		if !maxLifetimeSet {
			maxLifetime, maxLifetimeSet = defaultMaxLifetime, true
		}
	}

	if maxOpenSet {
		sqlDB.SetMaxOpenConns(maxOpen)
	}
	if maxIdleSet {
		sqlDB.SetMaxIdleConns(maxIdle)
	}
	if maxLifetimeSet {
		sqlDB.SetConnMaxLifetime(maxLifetime)
	}

	if maxOpenSet || maxIdleSet || maxLifetimeSet {
		slog.Info(
			"db pool configured",
			"maxOpenConns", sqlDB.Stats().MaxOpenConnections,
			"maxIdleConns", maxIdle,
			"connMaxLifetime", maxLifetime.String(),
		)
	}
}

func getOptionalIntEnv(key string) (value int, ok bool, err error) {
	raw := os.Getenv(key)
	if raw == "" {
		return 0, false, nil
	}
	v, err := strconv.Atoi(raw)
	if err != nil {
		return 0, false, fmt.Errorf("%s must be an integer: %w", key, err)
	}
	if v < 0 {
		return 0, false, fmt.Errorf("%s must be >= 0", key)
	}
	return v, true, nil
}

func getOptionalDurationEnv(key string) (value time.Duration, ok bool, err error) {
	raw := os.Getenv(key)
	if raw == "" {
		return 0, false, nil
	}
	d, err := time.ParseDuration(raw)
	if err != nil {
		return 0, false, fmt.Errorf("%s must be a Go duration (e.g. 30m, 1h): %w", key, err)
	}
	if d < 0 {
		return 0, false, fmt.Errorf("%s must be >= 0", key)
	}
	return d, true, nil
}
