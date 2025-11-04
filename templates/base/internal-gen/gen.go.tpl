// Package gen is used to work with different cryptographic algorithms
package gen

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return "", nil
	}

	hashedPassword := string(hash)

	return hashedPassword, nil
}

// GenerateRandomBytes returns securely generated random bytes.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
func GenerateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		return nil, err
	}

	return b, nil
}

// GenerateRandomString returns a URL-safe, base64 encoded
// securely generated random string.
func GenerateRandomString(s int) (string, error) {
	b, err := GenerateRandomBytes(s)
	return base64.URLEncoding.EncodeToString(b), err
}

// CompareHashAndPasswordBcrypt received hash and rawPassword then it checks both of them if strings match
// using bcrypt CompareHashAndPassword function, if strings match it returns nil, otherwise function returns error
func CompareHashAndPasswordBcrypt(hash, rawPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(rawPassword))
	if err != nil {
		return fmt.Errorf("incorrect password: %w", err)
	}

	return nil
}

