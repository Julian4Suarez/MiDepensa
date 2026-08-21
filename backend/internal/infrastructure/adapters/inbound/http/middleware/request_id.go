// Package middleware holds the cross-cutting HTTP concerns.
package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// RequestIDHeader carries the correlation id in and out of the service.
const RequestIDHeader = "X-Request-Id"

// contextKey is the gin context key under which the request id is stored.
const contextKey = "request_id"

// RequestID reuses an inbound correlation id or generates a new one.
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.GetHeader(RequestIDHeader)
		if id == "" {
			id = uuid.NewString()
		}
		c.Set(contextKey, id)
		c.Writer.Header().Set(RequestIDHeader, id)
		c.Next()
	}
}

// RequestIDFrom returns the correlation id attached to the request.
func RequestIDFrom(c *gin.Context) string {
	id, _ := c.Get(contextKey)
	value, _ := id.(string)
	return value
}
