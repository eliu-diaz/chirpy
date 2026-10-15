package auth

import (
	"net/http"
	"strings"
	"testing"
)

func TestAuth_GetBearerToken_ReturnsAnErrorWhenBearerTokenNotPresentInHeaders(t *testing.T) {
	// Arrange
	headers := http.Header{}

	// Act
	_, err := GetBearerToken(headers)

	// Assert
	if err == nil {
		t.Errorf("Expected GetBearerToken() to contain an error but it succeeded instead")
	}
}

func TestAuth_GetBearerToken_ReturnsTokenSegmentOnly(t *testing.T) {
	// Arrange
	headers := http.Header{}
	headerValue := "Bearer TOKEN_STRING"
	headers.Add("Authorization", headerValue)

	// Act
	token, err := GetBearerToken(headers)
	// Assert
	if err != nil {
		t.Fatalf("Expected GetBearerToken to return a token value but returned an error instead: %v", err)
	}

	splittedString := strings.Split(token, " ")
	if len(splittedString) > 1 {
		t.Errorf("Expected token to only contain the token part, instead, it contained various segments in it: %q", splittedString)
	}

	if strings.Split(headerValue, " ")[1] != token {
		t.Errorf("Received token didn't match the original token set as header value, token is %q and original header wasd %q", token, strings.Split(headerValue, " ")[1])
	}
}
