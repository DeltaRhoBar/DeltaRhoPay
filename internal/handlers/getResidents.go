package handlers

import (
	"deltapay/internal/services"
	"net/http"
	"encoding/json"
)

type GetResidentsHandler struct {
	db services.Database
}

func NewGetResidentHandler(db services.Database) *GetResidentsHandler {
	h := &GetResidentsHandler{db: db}	
	return h
}

func (h *GetResidentsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	residents, err := h.db.GetResidents();
	if err != nil {
		http.Error(w, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"residents": residents,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
