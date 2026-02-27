package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/joliverstrom-cmd/goServe/internal/auth"
	"github.com/joliverstrom-cmd/goServe/internal/database"
)

func (cfg *apiConfig) createPost(w http.ResponseWriter, req *http.Request) {

	type reqParameters struct {
		Body   string    `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}
	type returnVals struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Body      string    `json:"body"`
		UserID    uuid.UUID `json:"user_id"`
	}

	decoder := json.NewDecoder(req.Body)
	params := reqParameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 500, "Couldn't decode params", err)
		return
	}

	jwtString, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "No auth info in header", err)
		return
	}

	fmt.Println("validateJWE in posting")
	userID, err := auth.ValidateJWT(jwtString, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid token", err)
		return
	}

	if len(params.Body) > 140 {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

	cleanBody := stripString(params.Body)

	createdPost, err := cfg.db.CreatePost(req.Context(), database.CreatePostParams{
		Body:   cleanBody,
		UserID: uuid.NullUUID{UUID: userID, Valid: true},
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create db entry for post", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, returnVals{
		ID:        createdPost.ID,
		CreatedAt: createdPost.CreatedAt,
		UpdatedAt: createdPost.UpdatedAt,
		Body:      createdPost.Body,
		UserID:    createdPost.UserID.UUID,
	})

}

func stripString(body string) string {

	profanities := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}

	words := strings.Split(body, " ")
	lowerCaseWords := strings.Split(strings.ToLower(body), " ")

	for idx, word := range lowerCaseWords {
		if _, ok := profanities[word]; ok {
			words[idx] = "****"
		}
	}

	return strings.Join(words, " ")
}
