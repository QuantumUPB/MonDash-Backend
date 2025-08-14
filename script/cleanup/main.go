package main

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"

	"mondash-backend/logger"
	mongorepo "mondash-backend/repository/mongo"
)

func main() {
	_ = godotenv.Load()

	if err := logger.Init(); err != nil {
		log.Fatalf("failed to init logger: %v", err)
	}
	defer logger.Log.Sync()

	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		uri = "mongodb://mongodb:27017"
	}
	dbName := os.Getenv("MONGODB_DATABASE")
	if dbName == "" {
		dbName = "mondash"
	}

	db, err := mongorepo.Connect(uri, dbName)
	if err != nil {
		logger.Log.Fatalf("failed to connect to MongoDB: %v", err)
	}

	ctx := context.Background()

	collections := []string{
		"static_nodes",
		"node_history",
		"static_apps",
		"key_consumption",
		"map",
		"alerts_response",
		"auth_users",
		"device_keyrate",
	}

	for _, coll := range collections {
		if err := db.Collection(coll).Drop(ctx); err != nil {
			logger.Log.Warnw("failed to drop collection", "collection", coll, "err", err)
		}
	}

	logger.Log.Info("Database cleaned up successfully")
}
