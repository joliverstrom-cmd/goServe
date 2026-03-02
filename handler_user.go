package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/joliverstrom-cmd/goServe/internal/auth"
	"github.com/joliverstrom-cmd/goServe/internal/database"
)

func (cfg *apiConfig) createUser(w http.ResponseWriter, req *http.Request) {

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

	hashedPW, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create password hash", err)
	}

	createdUser, err := cfg.db.CreateUser(req.Context(), database.CreateUserParams{
		Email:          params.Email,
		HashedPassword: hashedPW,
	})
	if err != nil {
		respondWithError(w, 500, "Something went wrong", err)
		return
	}

	type returnVals struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Email     string    `json:"email"`
		ChirpyRed bool      `json:"is_chirpy_red"`
	}

	respondWithJSON(w, http.StatusCreated, returnVals{
		ID:        createdUser.ID,
		CreatedAt: createdUser.CreatedAt,
		UpdatedAt: createdUser.UpdatedAt,
		Email:     createdUser.Email,
		ChirpyRed: createdUser.IsChirpyRed.Bool,
	})

}

func (cfg *apiConfig) updateUser(w http.ResponseWriter, req *http.Request) {

	jwtString, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "No auth info in header", err)
		return
	}

	userID, err := auth.ValidateJWT(jwtString, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid token", err)
		return
	}

	type reqParameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	decoder := json.NewDecoder(req.Body)
	params := reqParameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 500, "Something went wrong", err)
		return
	}

	hashedPW, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create password hash", err)
	}

	updatedUser, err := cfg.db.UpdateUserDetails(req.Context(), database.UpdateUserDetailsParams{
		Email:          params.Email,
		HashedPassword: hashedPW,
		ID:             userID,
	})
	if err != nil {
		respondWithError(w, 500, "Something went wrong", err)
		return
	}

	type returnVals struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Email     string    `json:"email"`
		ChirpyRed bool      `json:"is_chirpy_red"`
	}

	respondWithJSON(w, http.StatusOK, returnVals{
		ID:        updatedUser.ID,
		CreatedAt: updatedUser.CreatedAt,
		UpdatedAt: updatedUser.UpdatedAt,
		Email:     updatedUser.Email,
		ChirpyRed: updatedUser.IsChirpyRed.Bool,
	})

}
