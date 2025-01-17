package main

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/vystepanenko/Chirpy/internal/auth"
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

	bar, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 401, "Unauthorized_bar")
		return
	}
	userId, err := auth.ValidateJWT(bar, cfg.secretKey)
	if err != nil {
		respondWithError(w, 401, "Unauthorized_user: "+err.Error())
		return
	}
	type requestBody struct {
		Body string `json:"body"`
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
		UserID: userId,
	}

	c, err := cfg.db.CreateChirps(context.Background(), chirpyParams)
	if err != nil {
		respondWithError(w, 400, err.Error())
		return
	}

	respondWithJSON(w, 201, Chirpy{
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
	auhorId := r.URL.Query().Get("author_id")
	var chirps []database.Chirp
	if auhorId != "" {
		userId, err := uuid.Parse(auhorId)
		if err != nil {
			respondWithError(w, 400, "Error getting chirps")
			return
		}
		chirps, err = cfg.db.GetAllChirpsByAuthor(context.Background(), userId)
		if err != nil {
			respondWithError(w, 400, "Error getting chirps")
			return
		}
	} else {
		var err error
		chirps, err = cfg.db.GetAllChirps(context.Background())
		if err != nil {
			respondWithError(w, 400, "Error getting chirps")
			return
		}
	}
	sort := r.URL.Query().Get("sort")
	if sort == "" {
		sort = "ASC"
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

	allChirps = sortChirpsByCreatedAt(allChirps, strings.ToUpper(sort))
	respondWithJSON(w, 200, allChirps)
}

func sortChirpsByCreatedAt(allChirps []Chirpy, order string) []Chirpy {
	sort.Slice(allChirps, func(i, j int) bool {
		if order == "DESC" {
			return allChirps[i].CreatedAt.After(allChirps[j].CreatedAt)
		}
		// Default to ASC
		return allChirps[i].CreatedAt.Before(allChirps[j].CreatedAt)
	})

	return allChirps
}

func (cfg *apiConfig) handlerGetChirp(w http.ResponseWriter, r *http.Request) {
	chirpId := r.PathValue("chirpId")
	if "" == chirpId {
		respondWithError(w, 422, "Please provide chirpId")
	}

	chirpIdParsed, err := uuid.Parse(chirpId)
	if err != nil {
		respondWithJSON(w, 404, "Chirp not found")
	}

	chirpDb, err := cfg.db.GetChirp(context.Background(), chirpIdParsed)
	if err != nil {
		respondWithJSON(w, 404, "Chirp not found")
	}

	respondWithJSON(w, 200, Chirpy{
		ID:        chirpDb.ID,
		UserId:    chirpDb.UserID,
		Body:      chirpDb.Body,
		CreatedAt: chirpDb.CreatedAt,
		UpdatedAt: chirpDb.UpdatedAt,
	})
}

func (cfg *apiConfig) handlerDeleteChirp(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	bar, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 401, "Unauthorized_bar")
		return
	}
	userId, err := auth.ValidateJWT(bar, cfg.secretKey)
	if err != nil {
		respondWithError(w, 401, "Unauthorized_user: "+err.Error())
		return
	}
	chirpId := r.PathValue("chirpId")
	if "" == chirpId {
		respondWithError(w, 422, "Please provide chirpId")
		return
	}

	chirpIdParsed, err := uuid.Parse(chirpId)
	if err != nil {
		respondWithError(w, 404, "Chirp not found")
		return
	}

	chirpDb, err := cfg.db.GetChirp(context.Background(), chirpIdParsed)
	if err != nil {
		respondWithError(w, 404, "Chirp not found")
	}

	if chirpDb.UserID != userId {
		respondWithError(w, 403, "Forbiden")
		return
	}

	dcParams := database.DeleteChirpParams{
		ID:     chirpDb.ID,
		UserID: userId,
	}
	err = cfg.db.DeleteChirp(context.Background(), dcParams)
	if err != nil {
		respondWithError(w, 400, "Something goes wrong: deleteChirp: "+err.Error())
	}

	w.WriteHeader(204)
}
