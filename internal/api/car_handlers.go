package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	db "github.com/Abhishekkapoor98/shimla-car-bazar/internal/db/sqlc"
)

func (server *Server) healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Shimla Car Bazar API is running with Clean Architecture! 🚗"))
}

func (server *Server) listCars(w http.ResponseWriter, r *http.Request) {
	cars, err := server.queries.ListCars(context.Background())
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to fetch cars: %v", err), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cars)
}

func (server *Server) createCar(w http.ResponseWriter, r *http.Request) {
	var req db.CreateCarParams
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf("Invalid input JSON: %v", err), http.StatusBadRequest)
		return
	}

	car, err := server.queries.CreateCar(context.Background(), req)
	if err != nil {
		log.Println("DB INSERT ERROR:", err)
		http.Error(w, fmt.Sprintf("Failed to save car: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(car)
}
