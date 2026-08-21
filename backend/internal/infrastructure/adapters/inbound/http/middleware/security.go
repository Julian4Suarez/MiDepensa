package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"midepensa/internal/config"
)

// SecurityHeaders sets the small set of headers that matter for a JSON API
// consumed by a browser SPA.
func SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("X-Content-Type-Options", "nosniff")
		c.Writer.Header().Set("X-Frame-Options", "DENY")
		c.Writer.Header().Set("Referrer-Policy", "no-referrer")
		c.Next()
	}
}

// CORS allows the configured frontend origins, or every origin when the
// allowlist is exactly "*" (local development only).
func CORS(cfg config.CORS) gin.HandlerFunc {
	allowAll := len(cfg.Origins) == 1 && cfg.Origins[0] == "*"
	allowed := make(map[string]struct{}, len(cfg.Origins))
	for _, origin := range cfg.Origins {
		allowed[origin] = struct{}{}
	}

	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		if origin != "" {
			if _, ok := allowed[origin]; ok || allowAll {
				c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
				c.Writer.Header().Set("Vary", "Origin")
				c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, OPTIONS")
				c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, "+RequestIDHeader)
				c.Writer.Header().Set("Access-Control-Max-Age", "600")
			}
		}

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}

// BodyLimit caps how much a client can send; the API only accepts small JSON.
func BodyLimit(maxBytes int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxBytes)
		c.Next()
	}
}
