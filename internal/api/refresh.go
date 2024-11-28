package api

import (
	"context"
	"database/sql"
	"net/http"
	"os"
	"time"

	"github.com/jmsMaupin1/chirpy/internal/auth"
	"github.com/jmsMaupin1/chirpy/internal/database"
)

func (cfg *ApiConfig) Refresh(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	type responseBody struct {
		Token string `json:"token"`
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		RespondWithError(w, 400, err.Error())
		return
	}

	refreshToken, err := cfg.DB.GetUserFromRefreshToken(context.Background(), token)
	if err != nil {
		RespondWithError(w, 400, err.Error())
		return
	}

	if (refreshToken.RevokedAt != sql.NullTime{}) || time.Now().After(refreshToken.ExpiresAt) {
		RespondWithError(w, 401, "You must log in again")
		return
	}

	newToken, err := auth.MakeJWT(refreshToken.ID, os.Getenv("SECRET"), time.Duration(1 * time.Hour))
	if err != nil {
		RespondWithError(w, 401, err.Error())
		return
	}

	RespondWithJson(w, 200, responseBody{
		Token: newToken,
	})
}

func (cfg *ApiConfig) Revoke(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	ctx := context.Background()

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		RespondWithError(w, 401, err.Error())
		return
	}

	user, err := cfg.DB.GetUserFromRefreshToken(ctx, token)
	if err != nil {
		RespondWithError(w, 401, err.Error())
		return
	}

	err = cfg.DB.RevokeRefreshToken(ctx, database.RevokeRefreshTokenParams{
		UserID: user.ID,
		RevokedAt: sql.NullTime{Time: time.Now() ,Valid: true},
		UpdatedAt: time.Now(),
	})
	if err != nil {
		RespondWithError(w, 401, err.Error())
	}

	w.WriteHeader(204)
}
