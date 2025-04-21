package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/niphawanphoopha/go-web-api/database"
	"github.com/niphawanphoopha/go-web-api/models"
)

// GetAllUsers returns all users (admin only)
func GetAllUsers(w http.ResponseWriter, r *http.Request) {
	var users []models.User
	
	// Get query parameters
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	
	// Set default values
	if limit <= 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}
	
	// Query the database
	query := database.DB.Model(&models.User{})
	
	// Apply pagination
	query = query.Limit(limit).Offset(offset)
	
	// Execute the query
	if err := query.Find(&users).Error; err != nil {
		http.Error(w, "Failed to fetch users", http.StatusInternalServerError)
		return
	}
	
	// Return the users
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(users)
}

// GetUserByID returns a user by ID (admin only)
func GetUserByID(w http.ResponseWriter, r *http.Request) {
	// Get the ID from the URL
	vars := mux.Vars(r)
	id := vars["id"]
	
	// Find the user in the database
	var user models.User
	if database.DB.First(&user, id).RecordNotFound() {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	
	// Return the user
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

// UpdateUser updates a user (admin only)
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	// Get the ID from the URL
	vars := mux.Vars(r)
	id := vars["id"]
	
	// Find the user in the database
	var user models.User
	if database.DB.First(&user, id).RecordNotFound() {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	
	// Parse request body
	var updatedUser models.User
	if err := json.NewDecoder(r.Body).Decode(&updatedUser); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	// Update the user
	user.FirstName = updatedUser.FirstName
	user.LastName = updatedUser.LastName
	user.Email = updatedUser.Email
	user.Role = updatedUser.Role
	
	// Save the updated user to the database
	if err := database.DB.Save(&user).Error; err != nil {
		http.Error(w, "Failed to update user", http.StatusInternalServerError)
		return
	}
	
	// Return the updated user
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

// DeleteUser deletes a user (admin only)
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	// Get the ID from the URL
	vars := mux.Vars(r)
	id := vars["id"]
	
	// Find the user in the database
	var user models.User
	if database.DB.First(&user, id).RecordNotFound() {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	
	// Delete the user from the database
	if err := database.DB.Delete(&user).Error; err != nil {
		http.Error(w, "Failed to delete user", http.StatusInternalServerError)
		return
	}
	
	// Return success message
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "User deleted"})
} 