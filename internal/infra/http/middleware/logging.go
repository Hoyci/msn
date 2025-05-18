package middleware

import (
	"log/slog"
	"msn/internal/config"
	"msn/internal/infra/logging"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
)

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg := config.GetConfig()
		start := time.Now()
		requestID := uuid.New().String()

		logger := logging.NewLogger(os.Stdout, cfg.Environment).With(
			slog.String("request_id", requestID),
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.String("remote_ip", r.RemoteAddr),
		)

		ctx := logging.ToContext(r.Context(), logger)
		r = r.WithContext(ctx)

		logger.InfoContext(ctx, "request_started")

		next.ServeHTTP(w, r)

		logger.InfoContext(ctx, "request_completed", slog.Duration("duration", time.Since(start)))
	})
}
