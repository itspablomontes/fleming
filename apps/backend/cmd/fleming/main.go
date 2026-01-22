package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/itspablomontes/fleming/api"
	"github.com/itspablomontes/fleming/api/internal/auth"
	"github.com/itspablomontes/fleming/api/internal/timeline"
)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})))

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		slog.Warn("DATABASE_URL not set")
	}

	port := os.Getenv("PORT")
	if port == "" {
		slog.Warn("PORT not set")
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

	if err := db.AutoMigrate(
		&auth.Challenge{},
		&timeline.TimelineEvent{},
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
