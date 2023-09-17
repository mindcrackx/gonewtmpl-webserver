package server

import (
	"fmt"
	"net/http"
)

func WithCacheControl(h http.Handler, maxAge int) http.Handler {
	cacheHeaderVal := fmt.Sprintf("public, max-age=%d", maxAge)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", cacheHeaderVal)
		h.ServeHTTP(w, r)
	})
}
