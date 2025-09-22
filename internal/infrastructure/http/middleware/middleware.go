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
	handler = AuthMiddleware(handler)
	handler = LoggingMiddleware(handler)
	handler = RecoveryMiddleware(handler)
	handler = CORS(ms.corsConfig)(handler)

	return handler
}

// ApplyToRouter applies middleware to a router
func (ms *MiddlewareStack) ApplyToRouter(router http.Handler) http.Handler {
	return ms.Apply(router)
}
