package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/niphawanphoopha/go-web-api/api"
	"github.com/niphawanphoopha/go-web-api/config"
	"github.com/niphawanphoopha/go-web-api/database"
	"github.com/niphawanphoopha/go-web-api/models"
)

// Item represents data about a record Item.
type Item struct {
	ID          string  `json:"id"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}

// Items slice to seed record Item data.
var items = []Item{
	{ID: "1", Title: "Item 1", Description: "This is item 1", Price: 19.99},
	{ID: "2", Title: "Item 2", Description: "This is item 2", Price: 29.99},
	{ID: "3", Title: "Item 3", Description: "This is item 3", Price: 39.99},
}

func main() {
	// Load configuration
	cfg := config.New()
	
	// Initialize database
	if err := database.Init(cfg); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.Close()
	
	// Auto-migrate models
	if err := database.AutoMigrate(&models.User{}, &models.Item{}); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}
	
	// Create a new server
	router := api.SetupRoutes(cfg)
	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Handler:      router,
		ReadTimeout:  time.Duration(cfg.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.WriteTimeout) * time.Second,
	}
	
	// Start the server in a goroutine
	go func() {
		log.Printf("Server starting on %s:%d...\n", cfg.Host, cfg.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()
	
	// Wait for interrupt signal to gracefully shut down the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	
	log.Println("Shutting down server...")
	
	// Create a deadline for server shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	// Shutdown the server
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
	
	log.Println("Server exiting")
}

func setupRoutes(router *gin.Engine) {
	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	// API group
	api := router.Group("/api")
	{
		// Items endpoints
		api.GET("/items", getItems)
		api.GET("/items/:id", getItemByID)
		api.POST("/items", createItem)
		api.PUT("/items/:id", updateItem)
		api.DELETE("/items/:id", deleteItem)
	}
}

// getItems responds with the list of all items as JSON.
func getItems(c *gin.Context) {
	c.JSON(http.StatusOK, items)
}

// getItemByID locates the item whose ID value matches the id
// parameter sent by the client, then returns that item as a response.
func getItemByID(c *gin.Context) {
	id := c.Param("id")

	// Loop over the list of items, looking for
	// an item whose ID matches the parameter.
	for _, item := range items {
		if item.ID == id {
			c.JSON(http.StatusOK, item)
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"message": "item not found"})
}

// createItem adds an item from JSON received in the request body.
func createItem(c *gin.Context) {
	var newItem Item

	// Call BindJSON to bind the received JSON to newItem.
	if err := c.BindJSON(&newItem); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Add the new item to the slice.
	items = append(items, newItem)
	c.JSON(http.StatusCreated, newItem)
}

// updateItem updates an item from JSON received in the request body.
func updateItem(c *gin.Context) {
	id := c.Param("id")
	var updatedItem Item

	if err := c.BindJSON(&updatedItem); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Loop through the items, looking for an item with matching ID
	for i, item := range items {
		if item.ID == id {
			updatedItem.ID = id
			items[i] = updatedItem
			c.JSON(http.StatusOK, updatedItem)
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"message": "item not found"})
}

// deleteItem removes an item from items slice.
func deleteItem(c *gin.Context) {
	id := c.Param("id")

	// Loop through the items, looking for an item with matching ID
	for i, item := range items {
		if item.ID == id {
			// Remove the item from the slice
			items = append(items[:i], items[i+1:]...)
			c.JSON(http.StatusOK, gin.H{"message": "Item deleted"})
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"message": "item not found"})
} 