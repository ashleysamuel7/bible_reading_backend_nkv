package main

import (
	"bible_reading_backend_nkv/database"
	"bible_reading_backend_nkv/models"
	"bible_reading_backend_nkv/server"
	"log"
	"os"
	_ "time/tzdata"
)

func main() {
	// Validate JWT_SECRET is set
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Println("WARNING: JWT_SECRET not set, using default secret. This should be changed in production!")
	}

	// Initialize database
	dbClient, err := database.NewDatabaseClient()
	if err != nil {
		log.Fatalf("failed to initialize Database Client: %v", err)
	}

	// Get the underlying GORM DB instance for migrations
	client, ok := dbClient.(*database.Client)
	if !ok {
		log.Fatalf("failed to get database client")
	}

	// Auto-migrate database tables
	log.Println("Running database migrations...")
	if err := client.DB.AutoMigrate(
		&models.User{},
		&models.UserFavoriteVerse{},
		&models.UserHighlightedVerse{},
		&models.UserLastRead{},
	); err != nil {
		log.Fatalf("failed to migrate database: %s", err)
	}
	log.Println("Database migrations completed successfully")

	// Create and start server
	serv := server.NewEchoServer(dbClient)
	if err := serv.Start(); err != nil {
		log.Fatal(err.Error())
	}
}