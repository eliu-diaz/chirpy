package main

import (
	"chirpy/internal/database"
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
	authorID := r.URL.Query().Get("author_id")
	var dbChirps []database.Chirp
	var dbGetChirpsError error

	if authorID != "" {
		authorUUID, parseError := uuid.Parse(authorID)
		if parseError != nil {
			respondWithError(w, http.StatusInternalServerError, "Could not parse author id from query parameters", parseError)
			return
		}

		dbChirps, dbGetChirpsError = cfg.db.RetrieveChirpsByAuthorID(r.Context(), authorUUID)
		if dbGetChirpsError != nil {
			respondWithError(w, http.StatusNotFound, fmt.Sprintf("Could not retrieve chirps from author %v", authorUUID), dbGetChirpsError)
			return
		}

	} else {
		dbChirps, dbGetChirpsError = cfg.db.RetrieveAllChirps(r.Context())
		if dbGetChirpsError != nil {
			respondWithError(w, http.StatusInternalServerError, "Could not retrieve chirps", dbGetChirpsError)
			return
		}
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
