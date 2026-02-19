package helpers

import (
	"crypto/sha256"
	"encoding/hex"
)

func HashPassword(password string) (string, error) {
	h := sha256.Sum256([]byte(password))
	return hex.EncodeToString(h[:]), nil
}

func CheckPasswordHash(password, hash string) bool {
	h, _ := HashPassword(password)
	return h == hash
}
