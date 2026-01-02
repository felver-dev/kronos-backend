package utils

import (
	"golang.org/x/crypto/bcrypt"
)

// HashPassword prend un mot de passe en clair et retourne son hash bcrypt
// Le coût par défaut de bcrypt est 10, ce qui est un bon équilibre sécurité/performance
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPasswordHash compare un mot de passe en clair avec un hash bcrypt
// Retourne true si le mot de passe correspond au hash, false sinon
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
