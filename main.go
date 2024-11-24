package main

import (
	"log"
	"net/http"

	"github.com/jmsMaupin1/chirpy/internal/api"
)

func main() {
	const port = "8080"
	mux := http.NewServeMux()

	cfg := api.ApiConfig{}

	srv := &http.Server{
		Handler: mux,
		Addr: ":" + port,
	}	

	mux.Handle("/", http.StripPrefix("/app", cfg.MiddlewareMetricsInc(http.FileServer(http.Dir(".")))))

	mux.HandleFunc("GET /metrics", cfg.MetricsHandler)
	mux.HandleFunc("POST /reset", cfg.MetricsReset)
	mux.HandleFunc("GET /test", func (w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Add("Content-Type", "text/plain; charset=utf-8")
		w.Header().Add("charset", "utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	log.Printf("Serving on port: %s\n", port)
	log.Fatal(srv.ListenAndServe())

}
