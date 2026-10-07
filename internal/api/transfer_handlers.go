package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	db "github.com/Abhishekkapoor98/shimla-car-bazar/internal/db/sqlc"
	"github.com/go-chi/chi/v5"
)

// Naya RC Transfer initiate karna (Jab koi gaadi kharidta hai)
func (server *Server) initiateTransfer(w http.ResponseWriter, r *http.Request) {
	var req db.InitiateRCTransferParams

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf("Invalid input JSON: %v", err), http.StatusBadRequest)
		return
	}

	transfer, err := server.queries.InitiateRCTransfer(context.Background(), req)
	if err != nil {
		log.Println("DB INSERT ERROR (Transfer):", err)
		http.Error(w, fmt.Sprintf("Failed to initiate transfer: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(transfer)
}

// Admin dashboard ke liye saare transfers track karna
func (server *Server) listTransfers(w http.ResponseWriter, r *http.Request) {
	transfers, err := server.queries.ListRCTransfers(context.Background())
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to fetch transfers: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(transfers)
}

// Admin RC transfer ka status update karega (e.g. 'in_rto_processing', 'completed')
func (server *Server) updateTransferStatus(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	transferID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid transfer ID", http.StatusBadRequest)
		return
	}

	var req struct {
		Status string `json:"status"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid input JSON", http.StatusBadRequest)
		return
	}

	updateArg := db.UpdateTransferStatusParams{
		ID:            transferID,
		CurrentStatus: req.Status,
	}

	updatedTransfer, err := server.queries.UpdateTransferStatus(context.Background(), updateArg)
	if err != nil {
		http.Error(w, "Failed to update transfer status", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedTransfer)
}
