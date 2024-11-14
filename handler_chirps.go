package main

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/vystepanenko/Chirpy/internal/database"
)

type Chirpy struct {
	ID        uuid.UUID `json:"id"`
	UserId    uuid.UUID `json:"user_id"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (cfg *apiConfig) handlerCreateChirp(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	type requestBody struct {
		Body   string    `json:"body"`
		UserId uuid.UUID `json:"user_id"`
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

	cleaned, err := validateChirpBody(params.Body)
	if err != nil {
		respondWithError(w, 400, err.Error())
		return
	}
	chirpyParams := database.CreateChirpsParams{
		Body:   cleaned,
		UserID: params.UserId,
	}

	c, err := cfg.db.CreateChirps(context.Background(), chirpyParams)
	if err != nil {
		respondWithError(w, 400, err.Error())
		return
	}

	respindWithJSON(w, 201, Chirpy{
		ID:        c.ID,
		UserId:    c.UserID,
		Body:      c.Body,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	},
	)
}

func validateChirpBody(body string) (string, error) {
	if len(body) > 140 {
		return "", errors.New("Chirp is to long")
	}

	return cleanBody(body), nil
}

func cleanBody(msg string) string {
	badWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}
	splitMsg := strings.Split(msg, " ")

	cleanedWords := make([]string, 0)
	for _, word := range splitMsg {
		if _, ok := badWords[strings.ToLower(word)]; ok {
			cleanedWords = append(cleanedWords, "****")
		} else {
			cleanedWords = append(cleanedWords, word)
		}
	}

	return strings.Join(cleanedWords, " ")
}

func (cfg *apiConfig) handlerGetAllChirps(w http.ResponseWriter, r *http.Request) {
	chirps, err := cfg.db.GetAllChirps(context.Background())
	if err != nil {
		respondWithError(w, 400, "Error getting chirps")
		return
	}

	allChirps := make([]Chirpy, 0)
	for _, v := range chirps {
		allChirps = append(allChirps, Chirpy{
			ID:        v.ID,
			UserId:    v.UserID,
			Body:      v.Body,
			CreatedAt: v.CreatedAt,
			UpdatedAt: v.UpdatedAt,
		})
	}

	respindWithJSON(w, 200, allChirps)
}

func (cfg *apiConfig) handlerGetChirp(w http.ResponseWriter, r *http.Request) {
	chirpId := r.PathValue("chirpId")
	if "" == chirpId {
		respondWithError(w, 422, "Please provide chirpId")
	}

	chirpIdParsed, err := uuid.Parse(chirpId)
	if err != nil {
		respindWithJSON(w, 404, "Chirp not found")
	}

	chirpDb, err := cfg.db.GetChirp(context.Background(), chirpIdParsed)
	if err != nil {
		respindWithJSON(w, 404, "Chirp not found")
	}

	respindWithJSON(w, 200, Chirpy{
		ID:        chirpDb.ID,
		UserId:    chirpDb.UserID,
		Body:      chirpDb.Body,
		CreatedAt: chirpDb.CreatedAt,
		UpdatedAt: chirpDb.UpdatedAt,
	})
}
