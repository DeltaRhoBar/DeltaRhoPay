package handlers

import (
	"deltapay/internal/services"
	"net/http"
)

type CheckoutHandler struct {
	db services.Database
}

func NewCheckoutHandler(db services.Database) *CheckoutHandler {
	h := &CheckoutHandler{db: db}	
	return h
}

func (h *CheckoutHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := h.db.CheckOut();
	if err != nil {
		http.Error(w, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
