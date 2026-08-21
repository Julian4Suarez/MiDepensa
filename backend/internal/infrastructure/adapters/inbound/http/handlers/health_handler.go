package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"midepensa/internal/application/ports"
	"midepensa/internal/config"
)

// checkTimeout bounds every dependency probe so /readyz can never hang.
const checkTimeout = 2 * time.Second

// HealthHandler serves the liveness and readiness probes.
type HealthHandler struct {
	checkers []ports.HealthChecker
}

// NewHealthHandler wires the probes to the dependencies to verify.
func NewHealthHandler(checkers ...ports.HealthChecker) *HealthHandler {
	return &HealthHandler{checkers: checkers}
}

// Healthz handles GET /healthz: the process is running.
func (h *HealthHandler) Healthz(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "ok",
		"commit":    config.Commit,
		"buildTime": config.BuildTime,
	})
}

// Readyz handles GET /readyz: every dependency answered in time.
func (h *HealthHandler) Readyz(c *gin.Context) {
	dependencies := make(map[string]string, len(h.checkers))
	ready := true

	for _, checker := range h.checkers {
		ctx, cancel := context.WithTimeout(c.Request.Context(), checkTimeout)
		err := checker.Check(ctx)
		cancel()

		if err != nil {
			ready = false
			dependencies[checker.Name()] = "down"
			continue
		}
		dependencies[checker.Name()] = "up"
	}

	status := http.StatusOK
	state := "ready"
	if !ready {
		status = http.StatusServiceUnavailable
		state = "not_ready"
	}
	c.JSON(status, gin.H{"status": state, "dependencies": dependencies})
}
