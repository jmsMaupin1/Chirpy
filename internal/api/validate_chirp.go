package api

import (
    "encoding/json"
    "fmt"
    "io"
    "net/http"
)

func Validate_json(w http.ResponseWriter, r *http.Request) {
    defer r.Body.Close()

    type requestBody struct {
	Body string `json:"body"`
    }

    type responseBody struct {
	Cleaned_body string `json:"cleaned_body"`
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


    if len(params.Body) == 0 {
	RespondWithError(w, 400, "You need to have more than 0 characters in the chirp")
	return
    } 

    if len(params.Body) > 140 {
	RespondWithError(w, 400, fmt.Sprintf("Chirps must be 140 characters, received chirp with %d characters", len(params.Body)))
	return
    }

    RespondWithJson(w, 200, responseBody{Cleaned_body: CleanChirp(params.Body)})
}
