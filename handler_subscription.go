package main

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/google/uuid"

	"github.com/vystepanenko/Chirpy/internal/auth"
	"github.com/vystepanenko/Chirpy/internal/database"
)

func (cfg *apiConfig) handlerSubscription(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil || apiKey != cfg.polkaKey {
		respondWithError(w, 401, "Unauthorized_bar")
		return
	}

	type requestBody struct {
		Event string `json:"event"`
		Data  struct {
			UserID string `json:"user_id"`
		} `json:"data"`
	}

	dat, err := io.ReadAll(r.Body)
	if err != nil {
		respondWithError(w, 400, "Something went wrong")
		return
	}

	params := requestBody{}
	err = json.Unmarshal(dat, &params)
	if err != nil {
		respondWithError(w, 400, "Something went wrong")
		return
	}

	if params.Event != "user.upgraded" {
		w.WriteHeader(204)
		return
	}

	userId, err := uuid.Parse(params.Data.UserID)
	if err != nil {
		respondWithError(w, 400, "Something went wrong")
		return
	}
	uiParams := database.UpdateChirpyRedParams{
		IsChirpyRed: true,
		ID:          userId,
	}
	_, err = cfg.db.UpdateChirpyRed(context.Background(), uiParams)
	if err != nil {
		w.WriteHeader(404)
		return
	}

	w.WriteHeader(204)
}
