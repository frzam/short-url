package models

import (
	"crypto/sha256"
	"fmt"
	"time"
)

func generateHash(originalURL string) string {
	// Generate current timestamp and appending it to url to create a unique hash even
	// for same url and return sha256 hash.
	currentTime := time.Now().Format("20060102150405")
	return fmt.Sprintf("%x", sha256.Sum256([]byte(originalURL+currentTime)))[:6]
}
