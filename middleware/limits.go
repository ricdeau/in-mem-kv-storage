package middleware

import (
	"github.com/ricdeau/in-mem-kv-storage/contracts"
	"github.com/ricdeau/in-mem-kv-storage/utils"
	"net/http"
)

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
