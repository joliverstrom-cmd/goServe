package main

import (
	"net/http"
	"time"

	"github.com/google/uuid"
)

type response struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) getPosts(w http.ResponseWriter, req *http.Request) {

	posts, err := cfg.db.GetAllChirps(req.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get posts from DB", err)
	}

	var myResponses []response

	for _, post := range posts {
		myResponses = append(myResponses, response{
			ID:        post.ID,
			CreatedAt: post.CreatedAt,
			UpdatedAt: post.UpdatedAt,
			Body:      post.Body,
			UserID:    post.UserID.UUID,
		})
	}

	respondWithJSON(w, http.StatusOK, myResponses)

}

func (cfg *apiConfig) getPost(w http.ResponseWriter, req *http.Request) {

	postID, err := uuid.Parse(req.PathValue("chirpID"))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't parse the provided uuid", err)
		return
	}

	post, err := cfg.db.GetOneChirp(req.Context(), postID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't find a post with that ID in DB", err)
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		ID:        post.ID,
		CreatedAt: post.CreatedAt,
		UpdatedAt: post.UpdatedAt,
		Body:      post.Body,
		UserID:    post.UserID.UUID,
	})

}
