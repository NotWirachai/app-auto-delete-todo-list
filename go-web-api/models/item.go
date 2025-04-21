package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

// Item represents data about a record Item.
type Item struct {
	gorm.Model
	ID          string    `json:"id" gorm:"primary_key"`
	Title       string    `json:"title" gorm:"not null"`
	Description string    `json:"description"`
	Price       float64   `json:"price" gorm:"not null"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// TableName specifies the table name for the Item model
func (Item) TableName() string {
	return "items"
}

// BeforeCreate is a GORM hook that runs before creating a new item
func (i *Item) BeforeCreate(scope *gorm.Scope) error {
	// You can add custom logic here, like validation
	return nil
}

// BeforeUpdate is a GORM hook that runs before updating an item
func (i *Item) BeforeUpdate(scope *gorm.Scope) error {
	// You can add custom logic here, like validation
	return nil
}

// Items is a slice of Item to simulate a database.
var Items = []Item{
	{ID: "1", Title: "Item 1", Description: "This is item 1", Price: 19.99},
	{ID: "2", Title: "Item 2", Description: "This is item 2", Price: 29.99},
	{ID: "3", Title: "Item 3", Description: "This is item 3", Price: 39.99},
} 