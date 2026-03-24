package ui

import (
	"crypto/rand"
	"encoding/hex"
)

// generateID generates a unique ID using crypto/rand
func generateID() string {
	b := make([]byte, 8)
	rand.Read(b)
	return hex.EncodeToString(b)
}
