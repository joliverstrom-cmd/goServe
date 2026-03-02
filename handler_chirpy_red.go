package main

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/joliverstrom-cmd/goServe/internal/auth"
)

const (
	event string = "user.upgraded"
)

func (cfg *apiConfig) setChirpyRedTrue(w http.ResponseWriter, req *http.Request) {

	givenKey, err := auth.GetAPIKey(req.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "No access", err)
	}

	if givenKey != cfg.polkaSecret {
		respondWithError(w, http.StatusUnauthorized, "No access", err)
	}

	type reqParameters struct {
		Event string `json:"event"`
		Data  struct {
			UserID string `json:"user_id"`
		} `json:"data"`
	}

	decoder := json.NewDecoder(req.Body)
	params := reqParameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 500, "Couldn't decode params", err)
		return
	}

	if params.Event != event {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	uuid, err := uuid.Parse(params.Data.UserID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Invalid UUID", err)
	}

	_, err = cfg.db.SetChirpyRedTrue(req.Context(), uuid)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Can't find a user with that UUID", err)
	}

	w.WriteHeader(http.StatusNoContent)

}
