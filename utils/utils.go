package utils

import (
	"encoding/json"
	"github.com/ricdeau/in-mem-kv-storage/contracts"
	"net/http"
)

const RequestId = "requestId"

func ErrorResponse(status int, error *contracts.Error, rw http.ResponseWriter) {
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(status)
	err := json.NewEncoder(rw).Encode(error)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
}
