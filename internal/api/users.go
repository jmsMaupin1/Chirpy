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

func (cfg *ApiConfig) AddUser(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	type requestBody struct {
		Email string `json:"email"`
	}

	type responseBody struct {
		ID string `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Email time.Time `json:"email"`
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

	user, err := cfg.DB.CreateUser(context.Background(), database.CreateUserParams{
		ID: uuid.New(),
		Email: params.Email,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})

	if err != nil {
		RespondWithError(w, 400, err.Error())
		return
	}

	RespondWithJson(w, 200, user)
}

func (cfg *ApiConfig) DeleteUser(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()	
}