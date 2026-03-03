package main

import (
	"fmt"
	"net/http"
	"slices"
	"time"

	"github.com/google/uuid"
	"github.com/joliverstrom-cmd/goServe/internal/database"
)

type response struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) getPosts(w http.ResponseWriter, req *http.Request) {

	author_id := req.URL.Query().Get("author_id")
	fmt.Printf("Here is the author-ID: %s\n", author_id)
	var posts []database.Post
	var err error

	if author_id != "" {

		pUUID, err := uuid.Parse(author_id)
		nullID := uuid.NullUUID{
			UUID:  pUUID,
			Valid: true,
		}

		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invlaid author-id", err)
		}
		posts, err = cfg.db.GetChirpsByAuthorID(req.Context(), nullID)
	} else {
		posts, err = cfg.db.GetAllChirps(req.Context())
	}
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

	sort_order := req.URL.Query().Get("sort")
	// Default from DB is "asc" = ascending, so we only need to sort if "desc"
	if sort_order == "desc" {
		slices.SortFunc(myResponses, func(a, b response) int {
			return b.CreatedAt.Compare(a.CreatedAt)
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
