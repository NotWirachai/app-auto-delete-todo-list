package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/niphawanphoopha/go-web-api/database"
	"github.com/niphawanphoopha/go-web-api/models"
)

// GetItems responds with the list of all items as JSON.
func GetItems(w http.ResponseWriter, r *http.Request) {
	var items []models.Item
	
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
	query := database.DB.Model(&models.Item{})
	
	// Apply pagination
	query = query.Limit(limit).Offset(offset)
	
	// Execute the query
	if err := query.Find(&items).Error; err != nil {
		http.Error(w, "Failed to fetch items", http.StatusInternalServerError)
		return
	}
	
	// Return the items
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(items)
}

// GetItemByID locates the item whose ID value matches the id
// parameter sent by the client, then returns that item as a response.
func GetItemByID(w http.ResponseWriter, r *http.Request) {
	// Get the ID from the URL
	vars := mux.Vars(r)
	id := vars["id"]
	
	// Find the item in the database
	var item models.Item
	if database.DB.Where("id = ?", id).First(&item).RecordNotFound() {
		http.Error(w, "Item not found", http.StatusNotFound)
		return
	}
	
	// Return the item
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(item)
}

// CreateItem adds an item from JSON received in the request body.
func CreateItem(w http.ResponseWriter, r *http.Request) {
	var item models.Item
	
	// Parse request body
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	// Validate request
	if item.Title == "" || item.Price <= 0 {
		http.Error(w, "Title and price are required", http.StatusBadRequest)
		return
	}
	
	// Save the item to the database
	if err := database.DB.Create(&item).Error; err != nil {
		http.Error(w, "Failed to create item", http.StatusInternalServerError)
		return
	}
	
	// Return the created item
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(item)
}

// UpdateItem updates an item from JSON received in the request body.
func UpdateItem(w http.ResponseWriter, r *http.Request) {
	// Get the ID from the URL
	vars := mux.Vars(r)
	id := vars["id"]
	
	// Find the item in the database
	var item models.Item
	if database.DB.Where("id = ?", id).First(&item).RecordNotFound() {
		http.Error(w, "Item not found", http.StatusNotFound)
		return
	}
	
	// Parse request body
	var updatedItem models.Item
	if err := json.NewDecoder(r.Body).Decode(&updatedItem); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	// Update the item
	item.Title = updatedItem.Title
	item.Description = updatedItem.Description
	item.Price = updatedItem.Price
	
	// Save the updated item to the database
	if err := database.DB.Save(&item).Error; err != nil {
		http.Error(w, "Failed to update item", http.StatusInternalServerError)
		return
	}
	
	// Return the updated item
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(item)
}

// DeleteItem removes an item from the database.
func DeleteItem(w http.ResponseWriter, r *http.Request) {
	// Get the ID from the URL
	vars := mux.Vars(r)
	id := vars["id"]
	
	// Find the item in the database
	var item models.Item
	if database.DB.Where("id = ?", id).First(&item).RecordNotFound() {
		http.Error(w, "Item not found", http.StatusNotFound)
		return
	}
	
	// Delete the item from the database
	if err := database.DB.Delete(&item).Error; err != nil {
		http.Error(w, "Failed to delete item", http.StatusInternalServerError)
		return
	}
	
	// Return success message
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Item deleted"})
} 