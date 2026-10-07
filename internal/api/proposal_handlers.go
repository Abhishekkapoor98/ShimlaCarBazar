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

func (server *Server) createProposal(w http.ResponseWriter, r *http.Request) {
	var req db.CreateSellProposalParams

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf("Invalid input JSON: %v", err), http.StatusBadRequest)
		return
	}

	proposal, err := server.queries.CreateSellProposal(context.Background(), req)
	if err != nil {
		log.Println("DB INSERT ERROR (Proposal):", err)
		http.Error(w, fmt.Sprintf("Failed to submit proposal: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(proposal)
}

func (server *Server) listProposals(w http.ResponseWriter, r *http.Request) {
	proposals, err := server.queries.ListSellProposals(context.Background())
	if err != nil {
		http.Error(w, "Failed to fetch proposals", http.StatusInternalServerError)
		return
	}

	role, _ := r.Context().Value(roleKey).(string)
	w.Header().Set("Content-Type", "application/json")

	// Agar Admin hai, toh saari requests bhej do
	if role == "admin" {
		json.NewEncoder(w).Encode(proposals)
		return
	}

	// Agar Customer hai, toh sirf uski requests filter karo
	userIDFloat, _ := r.Context().Value(userIDKey).(float64)

	var myProposals []db.SellProposal
	for _, p := range proposals {
		if fmt.Sprint(p.CustomerID) == fmt.Sprint(userIDFloat) {
			myProposals = append(myProposals, p)
		}
	}

	json.NewEncoder(w).Encode(myProposals)
}

func (server *Server) approveProposal(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	proposalID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid proposal ID", http.StatusBadRequest)
		return
	}

	proposal, err := server.queries.GetSellProposal(context.Background(), proposalID)
	if err != nil {
		http.Error(w, "Proposal not found", http.StatusNotFound)
		return
	}

	if proposal.Status == "approved" {
		http.Error(w, "This proposal is already approved", http.StatusBadRequest)
		return
	}

	// Proposal ka status 'approved' karna
	updateArg := db.UpdateSellProposalStatusParams{
		ID:     proposalID,
		Status: "approved",
	}
	updatedProposal, err := server.queries.UpdateSellProposalStatus(context.Background(), updateArg)
	if err != nil {
		http.Error(w, "Failed to approve proposal", http.StatusInternalServerError)
		return
	}

	// Inventory (cars table) mein automatic gaadi add karna, photos ke saath
	carArg := db.CreateCarParams{
		Title:         fmt.Sprintf("%s %s (2nd Hand)", proposal.CarBrand, proposal.CarModel),
		Brand:         proposal.CarBrand,
		Model:         proposal.CarModel,
		ConditionType: "old",
		Price:         proposal.ExpectedPrice,
		Status:        "available",
		ImageUrls:     proposal.ImageUrls, // Multi-photos database me ja rahi hain
	}

	newCar, err := server.queries.CreateCar(context.Background(), carArg)
	if err != nil {
		log.Println("Failed to add car to inventory:", err)
		http.Error(w, "Approved, but failed to create car entry", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"message":  "Proposal successfully approved and car added to inventory!",
		"car":      newCar,
		"proposal": updatedProposal,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
