package database

import (
	"to-do-list-app/models"
	"fmt"
)

func Migrate() {
	DB.AutoMigrate(&models.Todo{})
	fmt.Println("🛠️  Database Migrated")
}