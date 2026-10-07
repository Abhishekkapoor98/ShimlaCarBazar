package api

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

// Request context mein data save karne ke liye keys
type contextKey string

const (
	userIDKey contextKey = "user_id"
	roleKey   contextKey = "role"
)

// GUARD 1: Check karta hai ki User logged in hai (uske paas valid token hai)
func (server *Server) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 1. Header se token nikalna (Format: "Bearer eyJhbGciOi...")
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Unauthorized: Token missing", http.StatusUnauthorized)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Unauthorized: Invalid token format", http.StatusUnauthorized)
			return
		}
		tokenStr := parts[1]

		// 2. Token ko verify karna (wahi secret key use karke)
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			return jwtSecretKey, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Unauthorized: Invalid or expired token", http.StatusUnauthorized)
			return
		}

		// 3. Token se User ID aur Role nikal kar aage bhej dena
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "Unauthorized: Failed to parse claims", http.StatusUnauthorized)
			return
		}

		// User details ko request context me daal rahe hain taaki baaki APIs use kar sakein
		ctx := context.WithValue(r.Context(), userIDKey, claims["user_id"])
		ctx = context.WithValue(ctx, roleKey, claims["role"])

		// Next handler ko call karo with new context
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GUARD 2: Check karta hai ki User ek Admin hai ya nahi
func (server *Server) adminMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Context se role nikalna (jo pehle wale guard ne daala tha)
		role := r.Context().Value(roleKey)

		if role != "admin" {
			http.Error(w, "Forbidden: Only Admins can perform this action", http.StatusForbidden)
			return
		}

		// Agar admin hai, toh andar jaane do
		next.ServeHTTP(w, r)
	})
}
