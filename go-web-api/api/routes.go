package api

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/niphawanphoopha/go-web-api/config"
	"github.com/niphawanphoopha/go-web-api/handlers"
	"github.com/niphawanphoopha/go-web-api/middleware"
)

// SetupRoutes configures all the routes for our API
func SetupRoutes(cfg *config.Config) http.Handler {
	// Create a new router
	router := mux.NewRouter()
	
	// Add middleware
	router.Use(middleware.LoggingMiddleware)
	router.Use(middleware.CorsMiddleware())
	
	// Add config to context
	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), "config", cfg)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	})
	
	// Health check endpoint
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	}).Methods("GET")
	
	// API group
	api := router.PathPrefix("/api").Subrouter()
	
	// Auth routes (public)
	auth := api.PathPrefix("/auth").Subrouter()
	auth.HandleFunc("/register", handlers.Register).Methods("POST")
	auth.HandleFunc("/login", handlers.Login).Methods("POST")
	
	// Protected routes
	protected := api.PathPrefix("").Subrouter()
	protected.Use(middleware.AuthMiddleware(cfg))
	
	// User routes
	users := protected.PathPrefix("/users").Subrouter()
	users.HandleFunc("/me", handlers.GetCurrentUser).Methods("GET")
	
	// Admin routes
	admin := protected.PathPrefix("/admin").Subrouter()
	admin.Use(middleware.AdminMiddleware)
	
	// User management routes
	admin.HandleFunc("/users", handlers.GetAllUsers).Methods("GET")
	admin.HandleFunc("/users/{id}", handlers.GetUserByID).Methods("GET")
	admin.HandleFunc("/users/{id}", handlers.UpdateUser).Methods("PUT")
	admin.HandleFunc("/users/{id}", handlers.DeleteUser).Methods("DELETE")
	
	// Items routes
	items := protected.PathPrefix("/items").Subrouter()
	items.HandleFunc("", handlers.GetItems).Methods("GET")
	items.HandleFunc("/{id}", handlers.GetItemByID).Methods("GET")
	items.HandleFunc("", handlers.CreateItem).Methods("POST")
	items.HandleFunc("/{id}", handlers.UpdateItem).Methods("PUT")
	items.HandleFunc("/{id}", handlers.DeleteItem).Methods("DELETE")
	
	return router
} 