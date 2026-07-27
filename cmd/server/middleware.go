package main

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

func loggingMiddleware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(
		w http.ResponseWriter,
		r *http.Request,
	) {

		traceID := fmt.Sprintf(
			"TRACE-%d",
			time.Now().UnixNano(),
		)

		ctx := context.WithValue(
			r.Context(),
			"traceID",
			traceID,
		)

		r = r.WithContext(ctx)

		logger.Info(
			"Incoming Request",
			"traceID", traceID,
			"method", r.Method,
			"path", r.URL.Path,
		)

		next.ServeHTTP(w, r)
	})
}
