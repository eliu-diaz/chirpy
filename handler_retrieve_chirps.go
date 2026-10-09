package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) handleChirpRetrieval(w http.ResponseWriter, r *http.Request) {
	dbChirps, err := cfg.db.RetrieveAllChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not retrieve chirps", err)
		return
	}

	chirps := []chirp{}
	for _, dbChirp := range dbChirps {
		chirps = append(chirps, chirp{
			ID:        dbChirp.ID,
			CreatedAt: dbChirp.CreatedAt,
			UpdatedAt: dbChirp.UpdatedAt,
			Body:      dbChirp.Body,
			UserID:    dbChirp.UserID,
		})
	}

	respondWithJSON(w, http.StatusOK, chirps)
}

func (cfg *apiConfig) handleIndividualChirpRetrieval(w http.ResponseWriter, r *http.Request) {
	path := r.PathValue("chirpID")
	log.Printf("the passed path parameter is %s", path)

	// There's no way I can't convert this string uuid back to UUID type since it's created at the Database like that
	parsedUUID := uuid.MustParse(path)
	dbChirp, err := cfg.db.GetChirp(r.Context(), parsedUUID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, fmt.Sprintf("Could not find the chirp with UUID: %s", path), err)
		return
	}

	respondWithJSON(w, http.StatusOK, chirp{
		ID:        dbChirp.ID,
		CreatedAt: dbChirp.CreatedAt,
		UpdatedAt: dbChirp.UpdatedAt,
		Body:      dbChirp.Body,
		UserID:    dbChirp.UserID,
	})
}
