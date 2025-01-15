package api

import (
	"context"
	"encoding/json"
	"fmt"
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

func (cfg *ApiConfig) AddChirp() http.Handler {
	return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		type requestBody struct {
			Body string `json:"body"`
		}

		data, err := io.ReadAll(r.Body)
		if err != nil {
			RespondWithError(w, 400, err.Error())
			return
		}

		params := requestBody{}
		err = json.Unmarshal(data, &params)
		if err != nil{
			RespondWithError(w, 400, fmt.Sprintf("Unmarshalling error: %v", err.Error()))
			return
		}

		userID, err := uuid.Parse(r.Header.Get("user_id"))
		if err != nil {
			RespondWithError(w, 400, err.Error())
			return
		}

		if len(params.Body) == 0 || len(params.Body) > 140 {
			RespondWithError(w, 400, "Chirps must have at least 1 character and no more than 140 charactes")
			return
		}

		chirp, err := cfg.DB.CreateChirp(context.Background(), database.CreateChirpParams{
			ID: uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Body: CleanChirp(params.Body),
			UserID: userID,
		})

		if err != nil {
			RespondWithError(w, 400, fmt.Sprintf("Create Chirp Error: %v", err))
			return 
		}

		RespondWithJson(w, 201, Chirp(chirp))
	})
}

func (cfg *ApiConfig) DeleteChirp() http.Handler {
	return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		ctx := context.Background()

		idStr := r.Header.Get("user_id")
		uid, err := uuid.Parse(idStr)
		if err != nil {
			RespondWithError(w, 400, err.Error())
			return
		}

		chirpID, err := uuid.Parse(r.PathValue("id"))
		if err != nil {
			RespondWithError(w, 400, err.Error())
			return
		}

		chirp, err := cfg.DB.GetChirp(ctx, chirpID)
		if err != nil {
			RespondWithError(w, 404, err.Error())
			return
		}

		if chirp.UserID != uid {
			RespondWithError(w, 403, "Access Forbidden")
			return
		}

		err = cfg.DB.DeleteChirp(ctx, chirpID)
		if err != nil {
			RespondWithError(w, 400, err.Error())
		}

		w.WriteHeader(204)
	})
}

func (cfg *ApiConfig) GetChirp() http.Handler {
	return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		id, err := uuid.Parse(r.PathValue("id"))
		if err != nil {
			RespondWithError(w, 400, err.Error())
			return
		}

		c, err := cfg.DB.GetChirp(context.Background(), id)
		if err != nil {
			RespondWithError(w, 404, err.Error())
			return
		}

		RespondWithJson(w, 200, Chirp(c))
	})
}

func (cfg *ApiConfig) GetChirps() http.Handler {
	return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
	})
}

func (cfg *ApiConfig) GetAllChirps() http.Handler {
	return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		allChirps := []Chirp{}
		ctx := context.Background()

		aid := r.URL.Query().Get("author_id")
		sort := r.URL.Query().Get("sort")

		if sort == "" {
			sort = "desc"
		} 

		if aid != "" {
			uid, err := uuid.Parse(aid)
			if err != nil {
				RespondWithError(w, 400, err.Error())
				return
			}

			chirps, err := cfg.DB.GetChirpsByAuthor(ctx, database.GetChirpsByAuthorParams{
				UserID: uid,
				Sort: sort,
			})
			if err != nil {
				RespondWithError(w, 400, err.Error())
				return
			}

			for _, chirp := range chirps {
				allChirps = append(allChirps, Chirp(chirp))
			}
		} else {
			chirps, err := cfg.DB.GetAllChirps(context.Background(), sort)
			if err != nil {
				RespondWithError(w, 400, err.Error())
				return
			}

			for _, chirp := range chirps {
				allChirps = append(allChirps, Chirp(chirp))
			}
		}	

		RespondWithJson(w, 200, allChirps)
	})
}
