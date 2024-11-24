package api

import (
	"net/http"
	"sync/atomic"
	"fmt"
)

type ApiConfig struct {
	FileserveHits atomic.Int32
}

func (cfg *ApiConfig) MiddlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		cfg.FileserveHits.Add(1)
		next.ServeHTTP(w, req)
	})
}

func (cfg *ApiConfig) MetricsHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.Header().Add("charset", "utf-8")

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Hits: %d", cfg.FileserveHits.Load())))
}

func (cfg *ApiConfig) MetricsReset(w http.ResponseWriter, req *http.Request) {
	cfg.FileserveHits.Store(0)
	w.WriteHeader(http.StatusOK)
}
