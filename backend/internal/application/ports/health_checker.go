// Package ports declares outbound ports that are not persistence repositories.
package ports

import "context"

// HealthChecker reports whether a dependency is usable right now.
type HealthChecker interface {
	// Name identifies the dependency in the readiness payload.
	Name() string
	// Check returns nil when the dependency is healthy.
	Check(ctx context.Context) error
}
