package main

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
)

func main() {
	// Initialize logger
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()

	// Load configuration from environment variables
	cfg, err := LoadConfig()
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to load configuration")
	}

	// Create root context with logger
	ctx := context.Background()
	ctx = logger.WithContext(ctx)

	// Initialize GeoIP service
	geoService, err := NewGeoIPService(ctx, cfg)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to initialize GeoIP service")
	}
	defer geoService.Close()

	// Set up Mux router
	r := mux.NewRouter()

	// API routes
	r.HandleFunc("/health", geoService.HealthHandler).Methods("GET")
	r.HandleFunc("/lookup/{ip}", geoService.SingleIPHandler).Methods("GET")
	r.HandleFunc("/batch", geoService.BatchIPHandler).Methods("POST")

	// Middleware
	r.Use(loggingMiddleware(&logger))

	// Start server
	serverAddr := ":" + cfg.ServerPort
	logger.Info().Str("addr", serverAddr).Msg("Starting GeoIP API server")

	if err := http.ListenAndServe(serverAddr, r); err != nil {
		logger.Fatal().Err(err).Msg("Server failed to start")
	}
}

// Logging middleware
func loggingMiddleware(logger *zerolog.Logger) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// Wrap response writer to capture status code
			wrapped := &responseWriter{w, http.StatusOK}
			next.ServeHTTP(wrapped, r)

			// Log the request
			logger.Info().
				Str("method", r.Method).
				Str("url", r.URL.Path).
				Str("remote_addr", r.RemoteAddr).
				Int("status", wrapped.status).
				Str("user_agent", r.UserAgent()).
				Dur("duration", time.Since(start)).
				Msg("HTTP request")
		})
	}
}

// Custom response writer to capture status code
type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}
