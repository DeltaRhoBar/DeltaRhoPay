package handlers

import (
	"deltapay/internal/services"
	"net/http"
	"encoding/json"
)

type UpdateResidentHandler struct {
	db services.Database
}

type updatedResident struct {
	Id int
	R_floor int
	R_nr int
	Name string
	Telephone string
}

func NewUpdateResidentHandler(db services.Database) *UpdateResidentHandler {
	h := &UpdateResidentHandler{db: db}	
	return h
}

func (h *UpdateResidentHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var resident updatedResident
	if err := json.NewDecoder(r.Body).Decode(&resident); err != nil {
		http.Error(w, "Bad request: "+err.Error(), http.StatusBadRequest)
		return
	}

	err := h.db.UpdateResident(resident.Id, resident.R_floor, resident.R_nr, resident.Name, resident.Telephone)
	if err != nil {
		http.Error(w, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
