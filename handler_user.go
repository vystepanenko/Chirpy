package main

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"

	"github.com/vystepanenko/Chirpy/internal/auth"
	"github.com/vystepanenko/Chirpy/internal/database"
)

type User struct {
	ID          uuid.UUID `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Email       string    `json:"email"`
	IsChirpyRed bool      `json:"is_chirpy_red"`
}

func (cfg *apiConfig) handlerUserCreate(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	type requestBody struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	dat, err := io.ReadAll(r.Body)
	if err != nil {
		respondWithError(w, 400, "Something went wrong: body")
		return
	}

	params := requestBody{}
	err = json.Unmarshal(dat, &params)
	if err != nil {
		respondWithError(w, 400, "Something went wrong: unmarshal: "+err.Error())
		return
	}

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, 400, "Something went wrong: hashing password: "+err.Error())
		return
	}

	uParams := database.CreateUserParams{
		Email:          params.Email,
		HashedPassword: hashedPassword,
	}

	u, err := cfg.db.CreateUser(context.Background(), uParams)
	if err != nil {
		respondWithError(w, 400, err.Error())
		return
	}

	respondWithJSON(w, 201, User{
		ID:          u.ID,
		Email:       u.Email,
		CreatedAt:   u.CreatedAt,
		UpdatedAt:   u.UpdatedAt,
		IsChirpyRed: u.IsChirpyRed,
	},
	)
}

func (cfg *apiConfig) handlerUserLogin(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	type requestBody struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	type responseBody struct {
		User
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}

	dat, err := io.ReadAll(r.Body)
	if err != nil {
		respondWithError(w, 400, "Something went wrong: body")
		return
	}

	params := requestBody{}
	err = json.Unmarshal(dat, &params)
	if err != nil {
		respondWithError(w, 400, "Something went wrong: unmarshal: "+err.Error())
		return
	}

	user, err := cfg.db.GetUserByEmail(context.Background(), params.Email)
	if err != nil {
		respondWithError(w, 401, "Unauthorized")
		return
	}

	err = auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil {
		respondWithError(w, 401, "Unauthorized")
		return
	}

	accessToken, err := auth.MakeJWT(user.ID, cfg.secretKey, time.Hour)
	if err != nil {
		respondWithError(w, 401, "Unauthorized")
		return
	}

	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		respondWithError(w, 401, "Unauthorized: refresh")
		return
	}

	rtParams := database.CreateRefreshTokenParams{
		Token:  refreshToken,
		UserID: user.ID,
	}

	_, err = cfg.db.CreateRefreshToken(context.Background(), rtParams)
	if err != nil {
		respondWithError(w, 401, "Unauthorized: refresh_save: "+err.Error())
		return
	}

	respondWithJSON(w, 200, responseBody{
		User: User{
			ID:          user.ID,
			Email:       user.Email,
			CreatedAt:   user.CreatedAt,
			UpdatedAt:   user.UpdatedAt,
			IsChirpyRed: user.IsChirpyRed,
		},
		Token:        accessToken,
		RefreshToken: refreshToken,
	},
	)
}

func (cfg *apiConfig) handlerRefreshToken(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	bar, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 401, "Unauthorized: get bar")
		return
	}

	rt, err := cfg.db.GetRefreshToken(context.Background(), bar)
	if err != nil {
		respondWithError(w, 401, "Unauthorized: get from db: "+err.Error())
		return
	}

	accessToken, err := auth.MakeJWT(rt.UserID, cfg.secretKey, time.Hour)
	if err != nil {
		respondWithError(w, 401, "Unauthorized: make jwt")
		return
	}

	type responseBody struct {
		Token string `json:"token"`
	}
	respondWithJSON(w, 200, responseBody{
		Token: accessToken,
	})
}

func (cfg *apiConfig) handlerRefreshTokenRevoke(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	bar, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 401, "Unauthorized")
		return
	}

	rt, err := cfg.db.GetRefreshToken(context.Background(), bar)
	if err != nil {
		respondWithError(w, 401, "Unauthorized")
		return
	}

	err = cfg.db.RevokeRefreshToken(context.Background(), rt.Token)
	if err != nil {
		respondWithError(w, 401, "Unauthorized: revoke: "+err.Error())
		return
	}

	w.WriteHeader(204)
}

func (cfg *apiConfig) handlerUpdateUserInfo(w http.ResponseWriter, r *http.Request) {
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
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	dat, err := io.ReadAll(r.Body)
	if err != nil {
		respondWithError(w, 400, "Something went wrong: body")
		return
	}

	params := requestBody{}
	err = json.Unmarshal(dat, &params)
	if err != nil {
		respondWithError(w, 400, "Something went wrong: unmarshal: "+err.Error())
		return
	}

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, 400, "Something went wrong: hashing password: "+err.Error())
		return
	}

	uuParams := database.UpdateUserInfoParams{
		ID:             userId,
		Email:          params.Email,
		HashedPassword: hashedPassword,
	}

	u, err := cfg.db.UpdateUserInfo(context.Background(), uuParams)
	if err != nil {
		respondWithError(w, 401, "Unauthorized: updateInfo: "+err.Error())
	}

	respondWithJSON(w, 200, User{
		ID:          u.ID,
		Email:       u.Email,
		CreatedAt:   u.CreatedAt,
		UpdatedAt:   u.UpdatedAt,
		IsChirpyRed: u.IsChirpyRed,
	},
	)
}
