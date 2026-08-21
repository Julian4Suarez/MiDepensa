package middleware

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	"midepensa/internal/infrastructure/adapters/inbound/http/helpers"
)

// Recovery turns a panic into a 500 response instead of killing the process.
func Recovery(logger *slog.Logger) gin.HandlerFunc {
	return gin.CustomRecoveryWithWriter(nil, func(c *gin.Context, recovered any) {
		logger.Error("panic recovered",
			"path", c.Request.URL.Path,
			"request_id", RequestIDFrom(c),
			"panic", recovered,
		)
		helpers.RespondError(c, http.StatusInternalServerError, "internal_error", "something went wrong")
	})
}
