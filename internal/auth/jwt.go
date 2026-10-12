package auth

import (
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	claims := jwt.RegisteredClaims{
		Issuer:    "chirpy-access",
		Subject:   userID.String(),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(tokenSecret))
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	registeredClaims := jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(tokenString, &registeredClaims, func(t *jwt.Token) (any, error) {
		return []byte(tokenSecret), nil
	})
	if err != nil {
		return uuid.UUID{}, err
	}

	subject, err := token.Claims.GetSubject()
	if err != nil {
		log.Printf("An error occurred while trying to extract from claims, claims were %v", token.Claims)
		return uuid.UUID{}, err
	}

	userID, err := uuid.Parse(subject)
	if err != nil {
		log.Print("User claims didn't contain a properly formatted UUID")
		return uuid.UUID{}, err
	}
	return userID, nil
}
