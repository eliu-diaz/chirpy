package main

import (
	"chirpy/internal/auth"
	"net/http"
	"time"
)

func (cfg *apiConfig) handleRefresh(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Token string `json:"token"`
	}

	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get refresh token", err)
		return
	}
	dbRefreshToken, getTokenError := cfg.db.GetRefreshToken(r.Context(), refreshToken)
	if getTokenError != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't get refresh token", getTokenError)
		return
	}
	if dbRefreshToken.RevokedAt.Valid {
		respondWithError(w, http.StatusUnauthorized, "Refresh token was revoked", nil)
		return
	}
	// Has expired
	if dbRefreshToken.ExpiresAt.Before(time.Now()) {
		respondWithError(w, http.StatusUnauthorized, "Refresh token was expired", nil)
		return
	}

	newAccessToken, makeJWTError := auth.MakeJWT(dbRefreshToken.UserID, cfg.secret, time.Duration(3600)*time.Second)
	if makeJWTError != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not create token for user", makeJWTError)
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		Token: newAccessToken,
	})
}

func (cfg *apiConfig) handleRevoke(w http.ResponseWriter, r *http.Request) {
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
