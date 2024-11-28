package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/google/uuid"
	"github.com/jmsMaupin1/chirpy/internal/auth"
)

func (cfg *ApiConfig) WebhookPolka() http.Handler {
	return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		ctx := context.Background()

		type requestBody struct {
			Event string `json:"event"`
			Data struct {
				UserID string `json:"user_id"`
			} `json:"data"`
		}

		apiKey, err := auth.GetAPIKey(r.Header)
		if err != nil {
			RespondWithError(w, 401, err.Error())
			return 
		}

		if apiKey != os.Getenv("POLKA") {
			fmt.Println(apiKey)
			fmt.Println(os.Getenv("POLKA"))
			RespondWithError(w, 401, "Wrong api key")
			return
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

		if !strings.Contains(params.Event, "user.upgraded") {
			w.WriteHeader(204)
			return
		}

		uid, err := uuid.Parse(params.Data.UserID)
		if err != nil {
			RespondWithError(w, 400, err.Error())
			return
		}

		_, err = cfg.DB.GetUser(ctx, uid)
		if err != nil {
			RespondWithError(w, 404, "User not found")
			return
		}

		_, err = cfg.DB.SetUserChirpyRed(ctx, uid)
		if err != nil {
			RespondWithError(w, 400, err.Error())
		}

		w.WriteHeader(204)
	})
}
