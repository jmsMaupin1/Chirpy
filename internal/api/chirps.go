package api

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/jmsMaupin1/chirpy/internal/database"
)

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}


func (cfg *ApiConfig) AddChirp(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	type requestBody struct {
		Body string `json:"body"`
		User_id uuid.UUID `json:"user_id"`
	}

	data, err := io.ReadAll(r.Body)
	if err != nil {
		RespondWithError(w, 400, err.Error())
		return
	}

	params := requestBody{}
	err = json.Unmarshal(data, &params)
	if err != nil{
		RespondWithError(w, 400, err.Error())
		return
	}

	if len(params.Body) == 0 || len(params.Body) > 140 {
		RespondWithError(w, 400, "Chirps must have at least 1 character and no more than 140 charactes")
	}

	chirp, err := cfg.DB.CreateChirp(context.Background(), database.CreateChirpParams{
		ID: uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Body: CleanChirp(params.Body),
		UserID: params.User_id,
	})

	if err != nil {
		RespondWithError(w, 400, err.Error())
		return 
	}

	RespondWithJson(w, 201, Chirp(chirp))
}

func (cfg *ApiConfig) GetChirp(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		RespondWithError(w, 400, err.Error())
	}

	c, err := cfg.DB.GetChirp(context.Background(), id)
	if err != nil {
		RespondWithError(w, 404, err.Error())
	}

	RespondWithJson(w, 200, Chirp(c))
}

func (cfg *ApiConfig) GetAllChirps(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	allChirps := []Chirp{}

	chirps, err := cfg.DB.GetAllChirps(context.Background())
	if err != nil {
		RespondWithError(w, 400, err.Error())
		return
	}

	for _, chirp := range chirps {
		allChirps = append(allChirps, Chirp(chirp))
	}

	RespondWithJson(w, 200, allChirps)
}