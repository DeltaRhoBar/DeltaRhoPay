package services

import (
	"time"
	"crypto/rand"
	"encoding/hex"
	"os"
)


type Authenticator struct {
	tokens map[string]time.Time
	secret string
}

func NewAuthenticator() *Authenticator {
	secret := os.Getenv("SECRET")
	tokens := make(map[string]time.Time)
	return &Authenticator{secret: secret, tokens: tokens}
}

func (a *Authenticator) GenerateToken() (string, error) {
	exp := time.Now().Add(time.Hour)

	bytes := make([]byte, 16)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}

	tokenString := hex.EncodeToString(bytes)

	a.tokens[tokenString] = exp

	return tokenString, nil
}

func (a *Authenticator) Check(token string) bool {
	now := time.Now()
	for t, exp := range a.tokens {
		if now.After(exp) {
			delete(a.tokens, t)
		}
	}
	_, authorized := a.tokens[token]
	return authorized
}
