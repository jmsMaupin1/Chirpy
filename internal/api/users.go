package api

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/jmsMaupin1/chirpy/internal/auth"
	"github.com/jmsMaupin1/chirpy/internal/database"
)

type User struct {
	ID           uuid.UUID `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Email        string    `json:"email"`
	Token        string    `json:"token"`
	RefreshToken string    `json:"refresh_token"`
	IsChirpyRed  bool      `json:"is_chirpy_red"`
}

func (cfg *ApiConfig) AddUser() http.Handler {
	return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		type requestBody struct {
			Email string    `json:"email"`
			Password string `json:"password"`
		}

		type responseBody struct {
			ID string `json:"id"`
			CreatedAt time.Time `json:"created_at"`
			UpdatedAt time.Time `json:"updated_at"`
			Email time.Time     `json:"email"`
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

		hashedPass, err := auth.HashPassword(params.Password)
		if err != nil {
			RespondWithError(w, 400, err.Error())
		}

		user, err := cfg.DB.CreateUser(context.Background(), database.CreateUserParams{
			ID: uuid.New(),
			Email: params.Email,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			HashedPassword: hashedPass,
		})

		if err != nil {
			RespondWithError(w, 400, err.Error())
			return
		}

		RespondWithJson(w, 201, User{
			ID: user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email: user.Email,
		})
	})
}

func (cfg *ApiConfig) DeleteUsers(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()	

	if err := cfg.DB.DeleteUsers(context.Background()); err != nil {
		RespondWithError(w, 400, err.Error())
		return
	}

	RespondWithJson(w, 200, struct{Msg string}{
		Msg: "Success! Users deleted",
	})
}

func (cfg *ApiConfig) Login() http.Handler {
	return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		type requestBody struct {
			Email string		   `json:"email"`
			Password string		   `json:"password"`
		}

		ctx := context.Background()

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

		user, err := cfg.DB.GetUserByEmail(ctx, params.Email)
		if err != nil {
			RespondWithError(w, 400, err.Error())
			return
		}

		if err = auth.CheckPasswordHash(params.Password, user.HashedPassword); err != nil {
			RespondWithError(w, 401, err.Error())
			return
		}

		expiry_time := time.Duration(1 * time.Hour)

		accessToken, err := auth.MakeJWT(user.ID, cfg.Secret, expiry_time)
		if err != nil {
			RespondWithError(w, 400, err.Error())
			return
		}

		refreshToken, err := auth.MakeRefreshToken()
		if err != nil {
			RespondWithError(w, 400, err.Error())
		}

		cfg.DB.CreateRefreshToken(ctx, database.CreateRefreshTokenParams{
			Token: refreshToken,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			UserID: user.ID,
			ExpiresAt: time.Now().Add(time.Duration(60 * 24 * time.Hour)),
		})

		RespondWithJson(w, 200, User{
			ID: user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email: user.Email,
			Token: accessToken,
			RefreshToken: refreshToken,
			IsChirpyRed: user.IsChirpyRed,
		})
	})
}

func (cfg *ApiConfig) UpdateUser() http.Handler {
	return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		type requestBody struct {
			Email string `json:"email"`
			Password string `json:"password"`
		}

		bearer, err := auth.GetBearerToken(r.Header)
		if err != nil {
			RespondWithError(w, 401, err.Error())
			return
		}

		uid, err := auth.ValidateJWT(bearer, os.Getenv("SECRET"))
		if err != nil {
			RespondWithError(w, 401, err.Error())
			return
		}

		data, err := io.ReadAll(r.Body)
		if err != nil {
			RespondWithError(w, 401, err.Error())
			return
		}

		params := requestBody{}
		err = json.Unmarshal(data, &params)
		if err != nil{
			RespondWithError(w, 401, err.Error())
			return
		}

		hashed_password, err := auth.HashPassword(params.Password)
		if err != nil {
			RespondWithError(w, 401, err.Error())
			return
		}

		user, err := cfg.DB.UpdateUser(context.Background(), database.UpdateUserParams{
			ID: uid,
			UpdatedAt: time.Now(),
			HashedPassword: hashed_password,
			Email: params.Email,
		})

		if err != nil {
			RespondWithError(w, 401, err.Error())
			return
		}

		RespondWithJson(w, 200, User{
			ID: uid,
			UpdatedAt: time.Now(),
			CreatedAt: user.CreatedAt,
			Email: user.Email,
			Token: bearer,
		})
	})
}








