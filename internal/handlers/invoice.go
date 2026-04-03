package handlers

import (
	"bytes"
	"deltapay/internal/services"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

type InvoiceHandler struct {
	db services.Database
}

type invoiceData struct {
	Message string
	Ids []int
}

type message struct {
	Number  string `json:"number"`
	Message string `json:"message"`
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

	messages := []message{}
	for _, id := range invoiceData.Ids {
		messageData, err := h.db.GetDebt(id);
		if err != nil {
			log.Println(err)
			continue
		}
		messages = append(
			messages,
			message{
				Number: messageData.Telephone,
				Message: strings.NewReplacer(
					"{name}", messageData.Name,
					"{amount}", fmt.Sprintf("%.2f", float32(messageData.Amount) / 100),
					).Replace(invoiceData.Message),
			},
			)
	}

	jsonData, err := json.Marshal(messages)
	if err != nil {
		panic(err)
	}

	url := "http://localhost:4242/jobs"
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		http.Error(w, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	type JobResponse struct {
		JobID string `json:"jobId"`
	}

	var job JobResponse
	if err := json.Unmarshal(body, &job); err != nil {
		http.Error(w, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	type Job struct {
		LoginQRCode string `json:"loginQrCode"`
	}

	type NextJobResponse struct {
		ID  string `json:"id"`
		Job Job    `json:"job"`
	}

	var nextJob NextJobResponse

	for {
		time.Sleep(500 * time.Millisecond)
		client := &http.Client{}
		url = fmt.Sprintf("http://localhost:4242/jobs/%s", job.JobID)
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			http.Error(w, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
			return
		}
		resp, err = client.Do(req)
		if err != nil {
			http.Error(w, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()
		body, _ = io.ReadAll(resp.Body)
		if err := json.Unmarshal(body, &nextJob); err != nil {
			http.Error(w, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
			return
		}
		if nextJob.Job.LoginQRCode != "" {
			break
		}
	}
	response := map[string]any{
		"qr": nextJob.Job.LoginQRCode,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
