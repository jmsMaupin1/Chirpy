package main

import (
	"log"
	"net/http"
	"os"

	"github.com/jmsMaupin1/chirpy/internal/api"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	godotenv.Load()

	const port = "8080"
	mux := http.NewServeMux()

	cfg, err := api.New()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	srv := &http.Server{
		Handler: mux,
		Addr: ":" + port,
	}	

	mux.Handle("/app/", http.StripPrefix("/app", cfg.MiddlewareMetricsInc(http.FileServer(http.Dir(".")))))

	mux.HandleFunc("POST /api/chirps", cfg.AddChirp)
	mux.HandleFunc("POST /api/users", cfg.AddUser)
	mux.HandleFunc("GET /api/chirps/{id}", cfg.GetChirp)
	mux.HandleFunc("GET /api/healthz", api.Health)
	mux.HandleFunc("GET /api/chirps", cfg.GetAllChirps)

	mux.HandleFunc("GET /admin/metrics", cfg.MetricsHandler)
	mux.HandleFunc("POST /admin/reset", cfg.Reset)
	

	log.Printf("Serving on port: %s\n", port)
	log.Fatal(srv.ListenAndServe())

}
