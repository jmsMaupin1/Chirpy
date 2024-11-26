package api

import (
	"net/http"
	"html/template"
)

func (cfg *ApiConfig) MetricsHandler(w http.ResponseWriter, req *http.Request) {
		tmpl := template.Must(template.ParseFiles("layouts/admin/metrics/index.html"))
		tmpl.Execute(w, cfg.FileserveHits.Load())
}

func (cfg *ApiConfig) MetricsReset(w http.ResponseWriter, req *http.Request) {
		cfg.FileserveHits.Store(0)
		w.WriteHeader(http.StatusOK)
}
