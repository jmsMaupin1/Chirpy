package api

import (
	"net/http"
)

func Health(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.Header().Add("charset", "utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))	
}

func (cfg *ApiConfig) MiddlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		cfg.FileserveHits.Add(1)
		next.ServeHTTP(w, req)
	})
}
