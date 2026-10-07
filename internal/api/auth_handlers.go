package api

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	db "github.com/Abhishekkapoor98/shimla-car-bazar/internal/db/sqlc"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecretKey = []byte("shimla-car-bazar-super-secret-key-2026")

const OfficialAdminEmail = "admin@shimlacarbazar.com" // FIXED ADMIN EMAIL

func (server *Server) signupUser(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Role     string `json:"role"`
	}
	json.NewDecoder(r.Body).Decode(&req)

	// ADMIN VALIDATION: Agar koi Admin banna chahta hai toh email check karo
	if req.Role == "admin" && req.Email != OfficialAdminEmail {
		http.Error(w, "Access Denied: Only authorized person can be Admin.", http.StatusForbidden)
		return
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	arg := db.CreateUserParams{
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		Role:         req.Role,
	}

	_, err := server.queries.CreateUser(context.Background(), arg)
	if err != nil {
		http.Error(w, "Failed to create account. Email might already exist.", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "User created successfully"})
}

func (server *Server) loginUser(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	json.NewDecoder(r.Body).Decode(&req)

	user, err := server.queries.GetUserByEmail(context.Background(), req.Email)
	if err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	// DOUBLE CHECK: Agar DB me admin hai par email match nahi hoti (purane accounts ke liye)
	if user.Role == "admin" && user.Email != OfficialAdminEmail {
		http.Error(w, "Access Denied: You are not the official Admin.", http.StatusForbidden)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"role":    user.Role,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})
	tokenString, _ := token.SignedString(jwtSecretKey)

	json.NewEncoder(w).Encode(map[string]interface{}{
		"token": tokenString,
		"role":  user.Role,
		"email": user.Email,
	})
}
