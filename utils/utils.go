package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ricdeau/in-mem-kv-storage/contracts"
	"net/http"
)

// RequestID - key for storing unique request identifier
const RequestID = "requestId"

// ErrorResponse - creates error response with given status code, and writes it into response.
func ErrorResponse(status int, error *contracts.Error, rw http.ResponseWriter) {
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(status)
	err := json.NewEncoder(rw).Encode(error)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
	}
}

// RequestIDPrefix - tries to retrieve value with key 'RequestID' from context context.Context.
// Returns string like 'RequestID: 7e2e7a30-3465-4f21-9c99-51494a226596'
func RequestIDPrefix(ctx context.Context) string {
	requestID := ctx.Value(RequestID)
	return fmt.Sprintf("%s: %v", "RequestID", requestID)
}
