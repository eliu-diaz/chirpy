package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestAuth_MakeJWT_ReturnsAValidStringJWT(t *testing.T) {
	// Arrange
	userID := uuid.New()
	tokenSecret := "fido_loves_cats"
	expiresIn := time.Duration(4) * time.Hour

	// Act
	jwtResult, err := MakeJWT(userID, tokenSecret, expiresIn)
	// Assert
	if err != nil {
		t.Fatalf("Expected MakeJWT func to succeed while creating a new JWT token off of valid arguments, got an error instead: %v", err)
		return
	}

	if len(jwtResult) < 1 {
		t.Errorf("JWT was empty but it shoud've been created, got %q instead", jwtResult)
	}
}

func TestAuth_JWTRoundTrip_MakeJWTAndValidateIt(t *testing.T) {
	// Arrange
	userID := uuid.MustParse("38b2ece8-51ce-4d6b-af31-088719ed0dc7")
	tokenSecret := "fido_loves_cats"
	expiresIn := time.Duration(4) * time.Hour

	// Act
	jwtResult, err := MakeJWT(userID, tokenSecret, expiresIn)
	if err != nil {
		t.Fatalf("Expected MakeJWT func to succeed while creating a new JWT token off of valid arguments, got an error instead: %v", err)
	}

	receivedUUID, err := ValidateJWT(jwtResult, tokenSecret)
	if err != nil {
		t.Fatalf("Expected ValidateJWT to succeed, got an error instead: %v", err)
	}

	t.Logf("Received UUID from round trip is: %q and original UUID was: %q", receivedUUID, userID)
	if receivedUUID != userID {
		t.Errorf("Expected received UUID from JWT round trip to be the same, but receivedUUID was: %q and userID was: %q", receivedUUID, userID)
	}
}

// Test that validates expired tokens are rejected
func TestAuth_ValidateJWT_ExpiredTokenIsRejected(t *testing.T) {
	// Arrange
	userID := uuid.MustParse("38b2ece8-51ce-4d6b-af31-088719ed0dc7")
	tokenSecret := "fido_loves_cats"
	expiresIn := -(time.Duration(1) * time.Second)

	// Act
	jwtResult, err := MakeJWT(userID, tokenSecret, expiresIn)
	if err != nil {
		t.Fatalf("Expected MakeJWT func to succeed while creating a new JWT token off of valid arguments, got an error instead: %v", err)
	}

	_, validateErr := ValidateJWT(jwtResult, tokenSecret)
	if validateErr == nil {
		t.Errorf("Expected to get an error from ValidateJWT() since JWT was created as expired")
	}
}

func TestAuth_ValidateJWT_WrongSecretIsRejected(t *testing.T) {
	// Arrange
	userID := uuid.MustParse("38b2ece8-51ce-4d6b-af31-088719ed0dc7")
	tokenSecret := "fido_loves_cats"
	expiresIn := time.Duration(4) * time.Hour

	// Act
	jwtResult, err := MakeJWT(userID, tokenSecret, expiresIn)
	if err != nil {
		t.Fatalf("Expected MakeJWT func to succeed while creating a new JWT token off of valid arguments, got an error instead: %v", err)
	}

	wrongSecretKey := "fido_hates_cats"
	_, validateError := ValidateJWT(jwtResult, wrongSecretKey)

	if validateError == nil {
		t.Error("Expected ValidateJWT() to contain an error but it succeeded instead")
	}
}
