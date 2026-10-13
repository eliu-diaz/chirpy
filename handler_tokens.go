package main

import (
	"chirpy/internal/auth"
	"net/http"
)

func (cfg *apiConfig) revokeToken(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get refresh token", err)
		return
	}

	revokeTokenError := cfg.db.RevokeRefreshToken(r.Context(), refreshToken)
	if revokeTokenError != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't revoke refresh token", revokeTokenError)
		return
	}

	respondWithJSON(w, http.StatusNoContent, nil)
}
