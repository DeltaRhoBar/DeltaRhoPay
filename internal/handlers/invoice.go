package handlers

import (
	"deltapay/internal/services"
	"net/http"
	"encoding/json"
	"log"
)

type InvoiceHandler struct {
	db services.Database
}

type invoiceData struct {
	Message string
	Ids []int
}

func NewInvoiceHandler(db services.Database) *InvoiceHandler {
	h := &InvoiceHandler{db: db}	
	return h
}

func (h *InvoiceHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var invoiceData invoiceData 
	if err := json.NewDecoder(r.Body).Decode(&invoiceData); err != nil {
		log.Println(err)
		http.Error(w, "Bad request: "+err.Error(), http.StatusBadRequest)
		return
	}

	err := h.db.SetMessage(invoiceData.Message)
	if err != nil {
		http.Error(w, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
