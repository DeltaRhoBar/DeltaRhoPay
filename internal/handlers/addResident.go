package handlers

import (
	"deltapay/internal/services"
	"net/http"
	"encoding/json"
)

type AddResidentHandler struct {
	db services.Database
}

type newResident struct {
	R_floor int
	R_nr int
	Name string
}

func NewAddResidentHandler(db services.Database) *AddResidentHandler {
	h := &AddResidentHandler{db: db}	
	return h
}

func (h *AddResidentHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var resident newResident
	if err := json.NewDecoder(r.Body).Decode(&resident); err != nil {
		http.Error(w, "Bad request: "+err.Error(), http.StatusBadRequest)
		return
	}



	occupied, err := h.db.CheckOccupation(resident.R_floor, resident.R_nr)
	if err != nil {
		http.Error(w, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"occupied": occupied,
	}

	w.Header().Set("Content-Type", "application/json")

	if occupied {
		json.NewEncoder(w).Encode(response)
		return
	}

	err = h.db.AddResident(resident.R_floor, resident.R_nr, resident.Name)
	if err != nil {
		http.Error(w, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(response)
}
