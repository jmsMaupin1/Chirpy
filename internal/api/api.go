package api

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"os"
	"strings"
	"sync/atomic"

	"github.com/jmsMaupin1/chirpy/internal/database"
)

type ApiConfig struct {
	FileserveHits atomic.Int32
	DB database.Queries
}

func New() (*ApiConfig, error) {
	dbUrl := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		return nil, err
	}

	cfg := &ApiConfig{
		FileserveHits: atomic.Int32{},
		DB: *database.New(db),
	}

	return cfg, nil
	
}

func RespondWithJson(w http.ResponseWriter, status int, payload interface{}) error {
	res, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(res)

	return nil
}

func RespondWithError(w http.ResponseWriter, status int, error_msg string) error {
	return RespondWithJson(w, status, map[string]string{"error": error_msg})
}

func CleanChirp(s string) string {
	wordReplacements := map[string]string {
			"kerfuffle": "****",
			"sharbert": "****",
			"fornax": "****",
	}
	
	words := strings.Fields(s)

	for i, word := range words {
		if replacement, ok := wordReplacements[strings.ToLower(word)]; ok == true {
			words[i] = replacement
		}
	}

	return strings.Join(words, " ")
}
