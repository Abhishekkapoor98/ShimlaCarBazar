package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	db "github.com/Abhishekkapoor98/shimla-car-bazar/internal/db/sqlc"
	_ "github.com/lib/pq"

	// API aur DB package import kar rahe hain
	"github.com/Abhishekkapoor98/shimla-car-bazar/internal/api"
)

func main() {
	// Database connection setup
	dbSource := "postgresql://root:secretpassword@127.0.0.1:5433/car_bazar?sslmode=disable"
	conn, err := sql.Open("postgres", dbSource)
	if err != nil {
		log.Fatal("Cannot connect to db:", err)
	}

	if err = conn.Ping(); err != nil {
		log.Fatal("Database connection failed. Is Docker running?:", err)
	}
	defer conn.Close()

	// Queries aur API Server initialize karein
	queries := db.New(conn)
	server := api.NewServer(queries)

	// Server start karein
	port := ":8080"
	fmt.Printf("Shimla Car Bazar CDMS started on http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, server.Router()))
}
