package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/niphawanphoopha/go-web-api/config"
	"github.com/niphawanphoopha/go-web-api/database"
	"github.com/niphawanphoopha/go-web-api/middleware"
	"github.com/niphawanphoopha/go-web-api/models"
)

// RegisterRequest represents the request body for user registration
type RegisterRequest struct {
	Username  string `json:"username"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

// LoginRequest represents the request body for user login
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// AuthResponse represents the response for authentication endpoints
type AuthResponse struct {
	Token string      `json:"token"`
	User  models.User `json:"user"`
}

// Register handles user registration
func Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	
	// Parse request body
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	// Validate request
	if req.Username == "" || req.Email == "" || req.Password == "" {
		http.Error(w, "Username, email, and password are required", http.StatusBadRequest)
		return
	}
	
	// Check if username or email already exists
	var existingUser models.User
	if !database.DB.Where("username = ? OR email = ?", req.Username, req.Email).First(&existingUser).RecordNotFound() {
		http.Error(w, "Username or email already exists", http.StatusConflict)
		return
	}
	
	// Create new user
	user := models.User{
		Username:  req.Username,
		Email:     req.Email,
		Password:  req.Password, // Will be hashed by the BeforeCreate hook
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Role:      "user", // Default role
	}
	
	// Save user to database
	if err := database.DB.Create(&user).Error; err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}
	
	// Generate JWT token
	cfg := r.Context().Value("config").(*config.Config)
	token, err := middleware.GenerateToken(user.ID, user.Username, user.Role, cfg)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}
	
	// Return response
	response := AuthResponse{
		Token: token,
		User:  user,
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// Login handles user login
func Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	
	// Parse request body
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	// Validate request
	if req.Username == "" || req.Password == "" {
		http.Error(w, "Username and password are required", http.StatusBadRequest)
		return
	}
	
	// Find user by username
	var user models.User
	if database.DB.Where("username = ?", req.Username).First(&user).RecordNotFound() {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}
	
	// Check password
	if !user.CheckPassword(req.Password) {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}
	
	// Generate JWT token
	cfg := r.Context().Value("config").(*config.Config)
	token, err := middleware.GenerateToken(user.ID, user.Username, user.Role, cfg)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}
	
	// Return response
	response := AuthResponse{
		Token: token,
		User:  user,
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// GetCurrentUser returns the current user's information
func GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	// Get the claims from the context
	claims, ok := r.Context().Value("user").(*middleware.Claims)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	
	// Find user by ID
	var user models.User
	if database.DB.First(&user, claims.UserID).RecordNotFound() {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	
	// Return user information
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
} 