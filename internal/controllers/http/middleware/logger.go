package middleware

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

func Logger(next http.Handler, logger *zap.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t := time.Now()
		next.ServeHTTP(w, r)

		logger.Sugar().Infof("method=[%s] path=[%s] timeAnswer=%v", r.Method, r.URL.Path, time.Since(t))
	})
}
