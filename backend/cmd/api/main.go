// Command api is the MiDepensa HTTP service.
package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"midepensa/internal/application/services"
	"midepensa/internal/config"
	httpadapter "midepensa/internal/infrastructure/adapters/inbound/http"
	"midepensa/internal/infrastructure/adapters/outbound/persistence"
	"midepensa/internal/infrastructure/logging"
	"midepensa/internal/infrastructure/migrations"
)

// shutdownTimeout is how long in-flight requests get to finish on SIGTERM.
const shutdownTimeout = 10 * time.Second

func main() {
	// Absent in Docker and CI, where the environment is injected directly.
	_ = godotenv.Load()

	cfg := config.Load()
	logger := logging.New(cfg.Log)
	slog.SetDefault(logger)

	if cfg.App.AutoMigrate {
		if err := migrations.Run(cfg.Database.DSN()); err != nil {
			logger.Error("migrations failed", "error", err)
			os.Exit(1)
		}
		logger.Info("migrations applied")
	}

	ctx := context.Background()
	pool, err := persistence.NewPool(ctx, cfg.Database)
	if err != nil {
		logger.Error("database unavailable", "error", err)
		os.Exit(1)
	}
	defer pool.Close()

	pantryRepository := persistence.NewPostgresPantryRepository(pool)
	productRepository := persistence.NewPostgresProductRepository(pool)

	pantryService := services.NewPantryService(pantryRepository, productRepository)
	catalogService := services.NewCatalogService(productRepository)

	gin.SetMode(ginMode(cfg.Log.Level))
	router := httpadapter.SetupRouter(
		cfg,
		logger,
		pantryService,
		catalogService,
		persistence.NewPostgresHealthChecker(pool),
	)

	server := &http.Server{
		Addr:              cfg.App.Addr(),
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		logger.Info("server listening", "addr", server.Addr, "commit", config.Commit)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("server stopped", "error", err)
			os.Exit(1)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Error("graceful shutdown failed", "error", err)
	}
	logger.Info("server stopped")
}

func ginMode(logLevel string) string {
	if logLevel == "debug" {
		return gin.DebugMode
	}
	return gin.ReleaseMode
}
