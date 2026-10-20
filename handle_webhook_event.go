package main

import (
	"chirpy/internal/database"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handleWebhookEvent(w http.ResponseWriter, r *http.Request) {
	type data struct {
		UserID uuid.UUID `json:"user_id"`
	}

	type parameters struct {
		Event string `json:"event"`
		Data  data   `json:"data"`
	}

	params := parameters{}
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&params)

	switch params.Event {
	case "user.upgraded":
		// Returns the upgraded user, we don't need it now
		_, err := cfg.db.UpgradeToChirpyRed(r.Context(), database.UpgradeToChirpyRedParams{IsChirpyRed: true, ID: params.Data.UserID})
		if err != nil {
			respondWithError(w, http.StatusNotFound, "Could not upgrade user", err)
			return
		}

		respondWithJSON(w, http.StatusNoContent, nil)
	default:
		respondWithError(w, http.StatusNoContent, "Unsupported event", nil)
		return
	}

}
