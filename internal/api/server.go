package api

import (
	db "github.com/Abhishekkapoor98/shimla-car-bazar/internal/db/sqlc"
	"github.com/go-chi/chi/v5"
)

// Server struct database queries aur router ko hold karta hai
type Server struct {
	router  *chi.Mux
	queries *db.Queries
}

// NewServer ek naya server instance banata hai aur routes setup karta hai
func NewServer(queries *db.Queries) *Server {
	server := &Server{
		queries: queries,
		router:  chi.NewRouter(),
	}

	server.setupRoutes()
	return server
}

// Router return karta hai taaki main.go isey run kar sake
func (server *Server) Router() *chi.Mux {
	return server.router
}
