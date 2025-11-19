package handlers

import (
	"deltapay/internal/services"
	"net/http"
	"encoding/json"
	"log"
)

type PayHandler struct {
	db services.Database
}

type payData struct {
	Id int
}

func NewPayHandler(db services.Database) *PayHandler {
	h := &PayHandler{db: db}	
	return h
}

func (h *PayHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var payData payData
	if err := json.NewDecoder(r.Body).Decode(&payData); err != nil {
		log.Println(err)
		http.Error(w, "Bad request: "+err.Error(), http.StatusBadRequest)
		return
	}

	err := h.db.Pay(payData.Id)
	if err != nil {
		http.Error(w, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
