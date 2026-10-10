package auth

import (
	"testing"
)

// TestHashPassword passes a password argument to our test
// and ensures a proper hash is produced
func TestHashPassword(t *testing.T) {
	// Arrange
	password := "mypass"

	// Act
	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword received an error %v, it should've succeeded", err)
	}

	// Assert
	if hash == password {
		t.Errorf("Hash %q was exactly equal to the original password %q passed to hash", hash, password)
	}

	secondHash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("An error occurred while hashing password %q a second time, it should've succeeded but got an error instead %v", password, err)
	}

	if hash == secondHash {
		t.Errorf("Second Hash %q was equal to the original hash %q, which is unexpected for a hashing function", secondHash, hash)
	}
}
