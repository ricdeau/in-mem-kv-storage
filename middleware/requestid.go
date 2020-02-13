package middleware

import (
	"context"
	"github.com/google/uuid"
	"github.com/ricdeau/in-mem-kv-storage/logger"
	"github.com/ricdeau/in-mem-kv-storage/utils"
	"net/http"
)

// RequestIDMiddleware associate each request with unique ID and puts it in request's context.
func RequestIDMiddleware(next http.Handler) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		id := uuid.New()
		logger.Infof("RequestID: %v; Method: %s, path: %s, size: %d",
			id, r.Method, r.URL.Path, r.ContentLength)
		next.ServeHTTP(rw, r.WithContext(context.WithValue(r.Context(), utils.RequestID, id)))
	}
}
