package middleware

import (
    "net/http"
	"deltapay/internal/services"
)

func Auth(authenticator *services.Authenticator) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie("auth_token")
			if err != nil {
				if err == http.ErrNoCookie {
					// No cookie means user is not authenticated
					http.Redirect(w, r, "/login", http.StatusSeeOther)
					return
				}
				// Other error reading cookie
				http.Error(w, "Bad request", http.StatusBadRequest)
				return
			}
			token := cookie.Value
			ok := authenticator.Check(token)
			if !ok {
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}


