package auth

import (
	"fmt"
	"net/http"
	"strings"
)

func GetBearerToken(headers http.Header) (string, error) {
	token := headers.Get("Authorization")

	if token == "" {
		return token, fmt.Errorf("bearer token wasn't present in headers %v", headers)
	}

	splittedToken := strings.Split(token, " ")

	if len(splittedToken) > 1 {
		return splittedToken[1], nil
	} else {
		return "", fmt.Errorf("bearer token didn't have the right format, it contained: %q instead", token)
	}
}
