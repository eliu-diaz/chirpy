package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"slices"
	"strings"
	"time"

	"chirpy/internal/auth"
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
		token, err := auth.GetBearerToken(r.Header)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "No token provided", err)
			return
		}

		userID, err := auth.ValidateJWT(token, cfg.secret)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "Invalid token", err)
			return
		}

		chirp, err := extractChirp(r)
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
			Chirp{
				ID:        createdChirp.ID,
				CreatedAt: createdChirp.CreatedAt,
				UpdatedAt: createdChirp.UpdatedAt,
				Body:      createdChirp.Body,
				UserID:    createdChirp.UserID,
			},
		})
	})
}

func extractChirp(r *http.Request) (string, error) {
	type parameters struct {
		Body string `json:"body"`
	}
	var chirp string

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		return chirp, err
	}

	chirp = params.Body

	return chirp, nil
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
