package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/joliverstrom-cmd/goServe/internal/auth"
	"github.com/joliverstrom-cmd/goServe/internal/database"
)

func (cfg *apiConfig) refreshCheck(w http.ResponseWriter, req *http.Request) {

	refToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "No token in request: %v", err)
		return
	}

	dbToken, err := cfg.db.GetRefreshToken(req.Context(), refToken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid refresh token", err)
		return
	}

	if time.Now().Compare(dbToken.ExpiresAt) >= 0 {
		respondWithError(w, http.StatusUnauthorized, "Expired token", err)
		return
	}

	if dbToken.RevokedAt.Valid {
		respondWithError(w, http.StatusUnauthorized, "Expired token", err)
		return
	}

	fmt.Println("makeJWT in refreshCheck")
	userToken, err := auth.MakeJWT(dbToken.UserID.UUID, cfg.jwtSecret, time.Duration(time.Hour*1))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't generate user token", err)
		return
	}

	type returnVals struct {
		Token string `json:"token"`
	}

	respondWithJSON(w, http.StatusOK, returnVals{
		Token: userToken,
	})

}

func (cfg *apiConfig) revokeRefreshToken(w http.ResponseWriter, req *http.Request) {

	refToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "No token in request: %v", err)
		return
	}

	revokedToken, err := cfg.db.RevokeRefreshToken(req.Context(), database.RevokeRefreshTokenParams{
		RevokedAt: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
		Token: refToken,
	})
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't revoke token", err)
		return
	}
	type returned struct {
	}

	fmt.Printf("Token: %v, Updated at: %v, Revoked at: %v\n", revokedToken.Token, revokedToken.UpdatedAt, revokedToken.RevokedAt)

	respondWithJSON(w, http.StatusNoContent, returned{})

}
