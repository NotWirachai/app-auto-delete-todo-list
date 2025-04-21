package middleware

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/niphawanphoopha/go-web-api/config"
)

// Claims represents the JWT claims
type Claims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.StandardClaims
}

// GenerateToken generates a new JWT token for a user
func GenerateToken(userID uint, username, role string, cfg *config.Config) (string, error) {
	// Set expiration time
	expirationTime := time.Now().Add(time.Duration(cfg.JWTExpiry) * time.Minute)
	
	// Create claims
	claims := &Claims{
		UserID:   userID,
		Username: username,
		Role:     role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}
	
	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	
	// Sign token
	tokenString, err := token.SignedString([]byte(cfg.JWTSecret))
	if err != nil {
		return "", err
	}
	
	return tokenString, nil
}

// AuthMiddleware is a middleware that checks for a valid JWT token
func AuthMiddleware(cfg *config.Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get the Authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Authorization header is required", http.StatusUnauthorized)
				return
			}
			
			// Check if the header has the Bearer prefix
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				http.Error(w, "Authorization header format must be Bearer {token}", http.StatusUnauthorized)
				return
			}
			
			// Parse the token
			tokenString := parts[1]
			claims := &Claims{}
			
			token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
				return []byte(cfg.JWTSecret), nil
			})
			
			if err != nil {
				if err == jwt.ErrSignatureInvalid {
					http.Error(w, "Invalid token signature", http.StatusUnauthorized)
					return
				}
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}
			
			if !token.Valid {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}
			
			// Add the claims to the request context
			ctx := context.WithValue(r.Context(), "user", claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// AdminMiddleware is a middleware that checks if the user has the admin role
func AdminMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the claims from the context
		claims, ok := r.Context().Value("user").(*Claims)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		
		// Check if the user has the admin role
		if claims.Role != "admin" {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		
		next.ServeHTTP(w, r)
	})
} 