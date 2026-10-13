package main

import (
	"fmt"
	"log"

	"menu-tree-system/backend/config"
	"menu-tree-system/backend/handlers"
	"menu-tree-system/backend/models"
	"menu-tree-system/backend/repositories"
	"menu-tree-system/backend/routes"
	"menu-tree-system/backend/services"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// @title           Menu Tree System API
// @version         1.0
// @description     Hierarchical Menu Tree CRUD API with Infinite Nesting Support.
// @host            localhost:8080
// @BasePath        /
func main() {
	// 1. Load config
	cfg := config.LoadConfig()

	// 2. Connect to PostgreSQL
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// 3. Auto-migrate database schema
	log.Println("Running database migrations...")
	if err := db.AutoMigrate(&models.Menu{}); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	// 4. Initialize Dependency Injection (Repository -> Service -> Handler)
	menuRepo := repositories.NewMenuRepository(db)
	menuService := services.NewMenuService(menuRepo)
	menuHandler := handlers.NewMenuHandler(menuService)

	// 5. Setup Routes & Server
	r := routes.SetupRouter(cfg, menuHandler)

	serverAddr := ":" + cfg.ServerPort
	log.Printf("Server starting on port %s", cfg.ServerPort)
	if err := r.Run(serverAddr); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}