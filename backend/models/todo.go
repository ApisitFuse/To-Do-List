package models

import "gorm.io/gorm"

type Todo struct {
	gorm.Model
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
	DisplayOrder int `json:"displayOrder"`
}

type OrderObject struct {
	// gorm.Model
	OrId uint `json:"itemId"`
	NewIndex int `json:"newIndex"`
	OldIndex int `json:"oldIndex"`
}