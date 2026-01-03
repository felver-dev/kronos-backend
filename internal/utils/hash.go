package utils

import (
	"crypto/sha256"
	"encoding/hex"
)

// HashString génère un hash SHA256 d'une chaîne de caractères
// Utilisé pour hasher les tokens JWT avant de les stocker en base
func HashString(s string) string {
	hash := sha256.Sum256([]byte(s))
	return hex.EncodeToString(hash[:])
}

