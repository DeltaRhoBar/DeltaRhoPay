package services

import (
	"crypto/rand"
	"encoding/hex"
	"sync"
	"time"
)


type Authenticator struct {
	tokens map[string]time.Time
	mu sync.RWMutex
}

func NewAuthenticator() *Authenticator {
	tokens := make(map[string]time.Time)
	return &Authenticator{tokens: tokens}
}

func (a *Authenticator) GenerateToken() (string, error) {
	exp := time.Now().Add(time.Hour)

	bytes := make([]byte, 16)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}

	tokenString := hex.EncodeToString(bytes)

	a.mu.Lock()
	a.tokens[tokenString] = exp
	a.mu.Unlock()

	return tokenString, nil
}

func (a *Authenticator) Check(token string) bool {
	now := time.Now()

	a.mu.Lock()
	for t, exp := range a.tokens {
		if now.After(exp) {
			delete(a.tokens, t)
		}
	}
	_, authorized := a.tokens[token]
	a.mu.Unlock()

	return authorized
}
