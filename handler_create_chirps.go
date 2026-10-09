package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"slices"
	"strings"
	"time"

	"chirpy/internal/database"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handleChirpCreation() http.Handler {
	type Chirp struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Body      string    `json:"body"`
		UserID    uuid.UUID `json:"user_id"`
	}

	type response struct {
		Chirp
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		chirp, userID, err := extractChirpAndUserID(r)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters for chirp extraction", err)
			return
		}

		if !chirpIsValid(chirp) {
			respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Chirp \"%s\" is invalid", chirp), nil)
			return
		}

		createdChirp, err := cfg.db.CreateChirp(r.Context(), database.CreateChirpParams{Body: chirp, UserID: userID})
		if err != nil {
			formattedMessage := fmt.Sprintf("Couldn't create chirp %s with user id: %v, Error: %v", chirp, userID, err)
			respondWithError(w, http.StatusInternalServerError, formattedMessage, err)
			return
		}

		respondWithJSON(w, http.StatusCreated, response{
			Chirp: Chirp{
				ID:        createdChirp.ID,
				CreatedAt: createdChirp.CreatedAt,
				UpdatedAt: createdChirp.UpdatedAt,
				Body:      createdChirp.Body,
				UserID:    createdChirp.UserID,
			},
		})
	})
}

func extractChirpAndUserID(r *http.Request) (string, uuid.UUID, error) {
	type parameters struct {
		Body   string    `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}
	var chirp string
	var userID uuid.UUID

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		return chirp, userID, err
	}

	chirp = params.Body
	userID = params.UserID

	return chirp, userID, nil
}

func chirpIsValid(chirp string) bool {
	const maxChirpLength = 140

	return (len(chirp) <= maxChirpLength && len(chirp) > 0) && chirpDoesntHaveProfaneWord(chirp)
}

func chirpDoesntHaveProfaneWord(chirp string) bool {
	profaneWords := []string{"kerfuffle", "sharbert", "fornax"}
	splittedChirp := strings.Split(chirp, " ")

	for _, word := range splittedChirp {
		if slices.Contains(profaneWords, strings.ToLower(word)) {
			return false
		}
	}

	return true
}
