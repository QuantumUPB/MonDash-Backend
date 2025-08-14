package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"

	"mondash-backend/logger"
	mongorepo "mondash-backend/repository/mongo"
	"mondash-backend/routes"
)

func main() {
	// Load environment variables from .env if present
	_ = godotenv.Load()

	if err := logger.Init(); err != nil {
		log.Fatalf("failed to init logger: %v", err)
	}
	defer logger.Log.Sync()

	port := os.Getenv("PORT")
	if port == "" {
		port = "28080"
	}

	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://mongodb:27017"
	}
	dbName := os.Getenv("MONGODB_DATABASE")
	if dbName == "" {
		dbName = "mondash"
	}
	db, err := mongorepo.Connect(mongoURI, dbName)
	if err != nil {
		logger.Log.Fatalf("failed to connect to MongoDB: %v", err)
	}

	router := routes.NewRouter(db)

	logger.Log.Infof("Server started on :%s", port)
	if err := http.ListenAndServe(":"+port, router); err != nil {
		logger.Log.Fatal(err)
	}
}
