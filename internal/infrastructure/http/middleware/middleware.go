package middleware

import (
	"net/http"
	"time"
)

// MiddlewareStack holds all middleware
type MiddlewareStack struct {
	corsConfig  *CORSConfig
	rateLimiter *RateLimiter
}

// NewMiddlewareStack creates a new middleware stack
func NewMiddlewareStack() *MiddlewareStack {
	return &MiddlewareStack{
		corsConfig:  DefaultCORSConfig(),
		rateLimiter: NewRateLimiter(100, time.Minute), // 100 requests per minute
	}
}

// Apply applies all middleware to a handler
func (ms *MiddlewareStack) Apply(handler http.Handler) http.Handler {
	// Apply middleware in reverse order (last applied is first executed)
	handler = ms.rateLimiter.RateLimitMiddleware(handler)
	handler = LoggingMiddleware(handler)
	handler = RecoveryMiddleware(handler)
	handler = CORS(ms.corsConfig)(handler)

	return handler
}

// ApplyWithAuth applies middleware including authentication
func (ms *MiddlewareStack) ApplyWithAuth(handler http.Handler) http.Handler {
	handler = AuthMiddleware(handler)
	return ms.Apply(handler)
}

// ApplyWithLicense applies middleware including license validation
func (ms *MiddlewareStack) ApplyWithLicense(handler http.Handler) http.Handler {
	handler = LicenseMiddleware(handler)
	return ms.Apply(handler)
}

// ApplyPublic applies only basic middleware (no auth or license)
func (ms *MiddlewareStack) ApplyPublic(handler http.Handler) http.Handler {
	return ms.Apply(handler)
}

// ApplyToRouter applies middleware to a router
func (ms *MiddlewareStack) ApplyToRouter(router http.Handler) http.Handler {
	return ms.Apply(router)
}
