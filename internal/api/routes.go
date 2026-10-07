package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (server *Server) setupRoutes() {
	server.router.Use(middleware.Logger)
	server.router.Use(middleware.Recoverer)

	// ==========================================
	// 1. PUBLIC ROUTES (Koi bhi access kar sakta hai)
	// ==========================================
	server.router.Get("/health", server.healthCheck)
	server.router.Get("/cars", server.listCars) // Cars list public hai

	// Authentication APIs
	server.router.Post("/signup", server.signupUser)
	server.router.Post("/login", server.loginUser)

	// ==========================================
	// 2. PROTECTED ROUTES (Sirf Logged-in users ke liye)
	// ==========================================
	server.router.Group(func(r chi.Router) {
		// Guard 1: Token check karega
		r.Use(server.authMiddleware)

		// Upload API (Customer aur Admin dono photo upload kar sakte hain)
		r.Post("/upload", server.uploadImage)

		// Customer actions (Admin bhi access kar sakta hai)
		r.Post("/proposals", server.createProposal)
		r.Get("/proposals", server.listProposals)
		r.Post("/transfers", server.initiateTransfer)
		r.Get("/transfers", server.listTransfers)

		// ==========================================
		// 3. ADMIN ONLY ROUTES (Sirf Admin ke liye)
		// ==========================================
		r.Group(func(admin chi.Router) {
			// Guard 2: Role check karega (sirf 'admin' allowed)
			admin.Use(server.adminMiddleware)

			// Sirf Admin wale actions
			admin.Post("/cars", server.createCar)
			admin.Put("/proposals/{id}/approve", server.approveProposal)
			admin.Put("/transfers/{id}/status", server.updateTransferStatus)
		})
	})

	// ==========================================
	// 4. STATIC FILES & FRONTEND (404 Fix ke sath)
	// ==========================================

	// Photos serve karne ke liye
	fs := http.FileServer(http.Dir("uploads"))
	server.router.Handle("/uploads/*", http.StripPrefix("/uploads/", fs))

	// Frontend UI serve karne ke liye (Taki localhost:8080 par seedha website khule)
	server.router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "frontend/index.html")
	})
}
