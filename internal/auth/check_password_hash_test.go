package auth

import (
	"testing"
)

func TestCheckPasswordHash(t *testing.T) {
	// Act
	password := "fido123"

	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword produced an error while trying to hash %q. It should've produced a proper hash", password)
	}

	hashAndPasswordMatch, err := CheckPasswordHash(password, hash)
	if err != nil {
		t.Fatalf("CheckPasswordHash failed, it should've figured out if password %q and the hash %q matched but got error %v instead", password, hash, err)
	}
	if !hashAndPasswordMatch {
		t.Errorf("The password %q and the hash %v didn't match, CheckPasswordHash should've determined that these 2 matched", password, hash)
	}
}
