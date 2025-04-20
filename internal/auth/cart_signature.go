package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

func SignCartID(cartID string, secret []byte) string {
	mac := hmac.New(sha256.New, secret)
	mac.Write([]byte(cartID))
	return hex.EncodeToString(mac.Sum(nil))
}

func VerifyCartID(cartID, signature string, secret []byte) bool {
	expected := SignCartID(cartID, secret)
	return hmac.Equal([]byte(expected), []byte(signature))
}
