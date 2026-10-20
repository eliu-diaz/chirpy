package main

import (
	"chirpy/internal/auth"
	"chirpy/internal/database"
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handleChirpDeletion(w http.ResponseWriter, r *http.Request) {
	chirpIDString := r.PathValue("chirpID")
	chirpID, uuidParseError := uuid.Parse(chirpIDString)
	if uuidParseError != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't validate chirp", uuidParseError)
		return
	}

	token, authErr := auth.GetBearerToken(r.Header)
	if authErr != nil {
		respondWithError(w, http.StatusUnauthorized, "Could not find JWT", authErr)
		return
	}

	userID, jwtValidationError := auth.ValidateJWT(token, cfg.secret)
	if jwtValidationError != nil {
		respondWithError(w, http.StatusForbidden, "JWT was invalid", jwtValidationError)
		return
	}

	// We look for the chirp in the database so that we know who created it originally.
	dbChirp, getChirpError := cfg.db.GetChirp(r.Context(), chirpID)
	if getChirpError != nil {
		respondWithError(w, http.StatusNotFound, "Could not find chirp", getChirpError)
		return
	}

	if dbChirp.UserID != userID {
		respondWithError(w, http.StatusForbidden, "Author of the chirp is different than user trying to delete", nil)
		return
	}

	// Returns deleted chirp, unused
	deleteChirpError := cfg.db.DeleteChirp(r.Context(), database.DeleteChirpParams{ID: chirpID, UserID: userID})
	if deleteChirpError != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not delete chirp", deleteChirpError)
		return
	}

	respondWithJSON(w, http.StatusNoContent, nil)
}
