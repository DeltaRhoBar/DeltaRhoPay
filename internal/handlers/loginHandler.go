package handlers

import (
	"os"
	"net/http"
	"encoding/json"
	"deltapay/internal/services"
)

type loginData struct {
	Password string
}

type LoginHandler struct {
	key string
	auth *services.Authenticator
}

func NewLoginHandler(auth *services.Authenticator) *LoginHandler {
	key := os.Getenv("KEY")
	return &LoginHandler{key: key, auth: auth} 
}

func (h *LoginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var data loginData 
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "Bad request: "+err.Error(), http.StatusBadRequest)
		return
	}

	if (data.Password != h.key) {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	tokenString, err := h.auth.GenerateToken()
	if err != nil {
		http.Error(w, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    tokenString,
		Path:     "/",
		HttpOnly: true,        
		SameSite: http.SameSiteStrictMode,
		MaxAge:   3000,         
	})
	w.WriteHeader(http.StatusOK)
}
