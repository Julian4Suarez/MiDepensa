package helpers
// Package helpers translates domain errors into HTTP responses.
package helpers

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ErrorBody is the single error shape returned by every endpoint.
type ErrorBody struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type errorResponse struct {
	Error ErrorBody `json:"error"`
}

type mapping struct {
	status  int
	code    string
	message string
}

// registry is populated once during router setup, before any request is served.
var registry []struct {
	err     error
	mapping mapping
}

// RegisterDomainError declares how a sentinel error is rendered over HTTP.
func RegisterDomainError(err error, status int, code, message string) {
	registry = append(registry, struct {
		err     error
		mapping mapping
	}{err: err, mapping: mapping{status: status, code: code, message: message}})
}

// Respond writes the mapped response for err, defaulting to 500 for unknown
// errors so internal details never leak to the client.
func Respond(c *gin.Context, err error) {
	for _, entry := range registry {
		if errors.Is(err, entry.err) {
			c.AbortWithStatusJSON(entry.mapping.status, errorResponse{Error: ErrorBody{
				Code:    entry.mapping.code,
				Message: entry.mapping.message,
			}})
			return
		}
	}

	slog.Error("unhandled error", "path", c.FullPath(), "error", err)
	RespondError(c, http.StatusInternalServerError, "internal_error", "something went wrong")
}

// RespondError writes an explicit error response.
func RespondError(c *gin.Context, status int, code, message string) {
	c.AbortWithStatusJSON(status, errorResponse{Error: ErrorBody{Code: code, Message: message}})
}
