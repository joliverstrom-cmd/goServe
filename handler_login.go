package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/joliverstrom-cmd/goServe/internal/auth"
	"github.com/joliverstrom-cmd/goServe/internal/database"
)

func (cfg *apiConfig) login(w http.ResponseWriter, req *http.Request) {

	type reqParameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	decoder := json.NewDecoder(req.Body)
	params := reqParameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 500, "Something went wrong", err)
		return
	}

	defaultExpiryTimeSeconds := 60 * 60

	userDetails, err := cfg.db.GetUserByMail(req.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	check, err := auth.CheckPasswordHash(params.Password, userDetails.HashedPassword)
	if !check || err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}
	fmt.Println("makeJWT in login")
	userToken, err := auth.MakeJWT(userDetails.ID, cfg.jwtSecret, time.Duration(time.Second*time.Duration(defaultExpiryTimeSeconds)))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't generate user token", err)
		return
	}

	refreshToken := auth.MakeRefreshToken()

	dbRefreshToken, err := cfg.db.CreateRefreshToken(req.Context(), database.CreateRefreshTokenParams{
		Token: refreshToken,
		UserID: uuid.NullUUID{
			UUID:  userDetails.ID,
			Valid: true},
		ExpiresAt: time.Now().Add(time.Hour * 24 * 60),
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't add refresh token", err)
		return
	}

	type returnVals struct {
		ID           uuid.UUID `json:"id"`
		CreatedAt    time.Time `json:"created_at"`
		UpdatedAt    time.Time `json:"updated_at"`
		Email        string    `json:"email"`
		UserToken    string    `json:"token"`
		RefreshToken string    `json:"refresh_token"`
	}

	respondWithJSON(w, http.StatusOK, returnVals{
		ID:           userDetails.ID,
		CreatedAt:    userDetails.CreatedAt,
		UpdatedAt:    userDetails.CreatedAt,
		Email:        userDetails.Email,
		UserToken:    userToken,
		RefreshToken: dbRefreshToken.Token,
	})

}
