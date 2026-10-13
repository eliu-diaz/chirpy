package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"chirpy/internal/auth"
	"chirpy/internal/database"

	"github.com/google/uuid"
)

func (cfg *apiConfig) createUser() http.Handler {
	type User struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Email     string    `json:"email"`
	}

	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	type response struct {
		User
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		params := parameters{}
		err := decoder.Decode(&params)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
			return
		}

		hashedPassword, err := auth.HashPassword(params.Password)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Couldn't hash password %q", params.Password), err)
		}

		createUserParams := database.CreateUserParams{Email: params.Email, HashedPassword: hashedPassword}
		user, err := cfg.db.CreateUser(r.Context(), createUserParams)
		if err != nil {
			formattedMessage := fmt.Sprintf("Couldn't create user with email: %s", params.Email)
			respondWithError(w, http.StatusInternalServerError, formattedMessage, err)
			return
		}

		respondWithJSON(w, http.StatusCreated, response{
			User{
				ID:        user.ID,
				CreatedAt: user.CreatedAt,
				UpdatedAt: user.UpdatedAt,
				Email:     user.Email,
			},
		})
	})
}

func (cfg *apiConfig) loginUser(w http.ResponseWriter, r *http.Request) {
	type User struct {
		ID           uuid.UUID `json:"id"`
		CreatedAt    time.Time `json:"created_at"`
		UpdatedAt    time.Time `json:"updated_at"`
		Email        string    `json:"email"`
		AccessToken  string    `json:"token"`
		RefreshToken string    `json:"refresh_token"`
	}

	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	type response struct {
		User
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	dbUser, err := cfg.db.GetUser(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, fmt.Sprintf("Incorrect email %q", params.Email), err)
		return
	}

	ok, err := auth.CheckPasswordHash(params.Password, dbUser.HashedPassword)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, fmt.Sprintf("Incorrect password with error %v", err), err)
		return
	}

	// 1 hour expiration
	expiresIn := time.Duration(3600) * time.Second
	accessToken, makeJWTError := auth.MakeJWT(dbUser.ID, cfg.secret, expiresIn)
	if makeJWTError != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not create token for user", makeJWTError)
		return
	}

	if len(accessToken) == 0 {
		respondWithError(w, http.StatusInternalServerError, "User token was empty", nil)
		return
	}

	refreshToken := auth.MakeRefreshToken()
	refreshTokenExpiration := time.Now().Add((time.Duration(24) * time.Hour) * 60)
	createdRefreshToken, err := cfg.db.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token:     refreshToken,
		UserID:    dbUser.ID,
		ExpiresAt: refreshTokenExpiration,
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create refresh token", err)
	}

	if ok {
		respondWithJSON(w, http.StatusOK, response{
			User{
				ID:           dbUser.ID,
				CreatedAt:    dbUser.CreatedAt,
				UpdatedAt:    dbUser.UpdatedAt,
				Email:        dbUser.Email,
				AccessToken:  accessToken,
				RefreshToken: createdRefreshToken.Token,
			},
		})
	} else {
		respondWithError(w, http.StatusUnauthorized, "Incorrect password", nil)
	}
}
