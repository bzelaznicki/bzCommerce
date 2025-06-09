package auth

import (
	"fmt"
	"net/http"
	"strings"
)

func GetBearerToken(headers http.Header) (string, error) {
	header := strings.Split(headers.Get("Authorization"), " ")

	if len(header) < 2 {
		return "", fmt.Errorf("no valid auth header")
	}

	authToken := strings.TrimSpace(header[1])

	if len(authToken) == 0 {
		return "", fmt.Errorf("auth token missing")
	}

	return authToken, nil
}
