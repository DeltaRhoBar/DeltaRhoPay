package handlers

import (
	"deltapay/internal/services"
	"net/http"
	"encoding/json"
)

type ForceAddResidentHandler struct {
	db services.Database
}

func NewForceAddResidentHandler(db services.Database) *ForceAddResidentHandler {
	h := &ForceAddResidentHandler{db: db}	
	return h
}

func (h *ForceAddResidentHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var resident newResident
	if err := json.NewDecoder(r.Body).Decode(&resident); err != nil {
		http.Error(w, "Bad request: "+err.Error(), http.StatusBadRequest)
		return
	}

	err := h.db.AddResidentReplace(resident.R_floor, resident.R_nr, resident.Name, resident.Telephone)
	if err != nil {
		http.Error(w, "Internal Server Error "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.WriteHeader(http.StatusOK)
}
