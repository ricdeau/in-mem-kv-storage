package middleware

import (
	"github.com/ricdeau/in-mem-kv-storage/contracts"
	"github.com/ricdeau/in-mem-kv-storage/utils"
	"net/http"
)

// LimitsMiddleware creates middleware that filters incoming requests by key and value sizes.
// prefix - part of the uri before the key part.
// maxKey, maxValue - maximum key and value sizes in bytes, zero or below values means unlimited sizes.
// If key or value exceeds given limits writes error in json format to response and returns 414 status code.
func LimitsMiddleware(next http.Handler, prefix string, maxKey, maxValue int) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		keySize := len(r.URL.Path) - len(prefix)
		if maxKey > 0 && keySize > maxKey {
			utils.ErrorResponse(http.StatusRequestURITooLong, contracts.KeyToLongError(keySize, maxKey), rw)
			return
		}
		valSize := int(r.ContentLength)
		if maxValue > 0 && valSize > maxValue {
			utils.ErrorResponse(http.StatusRequestURITooLong, contracts.ValueToLongError(valSize, maxValue), rw)
			return
		}
		next.ServeHTTP(rw, r)
	}
}
