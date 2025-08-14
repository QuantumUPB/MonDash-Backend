package main

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"

	"mondash-backend/logger"
	"mondash-backend/repository/inmemory"
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

	// clean existing data
	db.Collection("static_nodes").Drop(ctx)
	db.Collection("node_history").Drop(ctx)
	db.Collection("static_apps").Drop(ctx)
	db.Collection("key_consumption").Drop(ctx)
	db.Collection("map").Drop(ctx)
	db.Collection("alerts_response").Drop(ctx)
	db.Collection("auth_users").Drop(ctx)

	// populate static node information
	if nodes := inmemory.DefaultNodeData(); len(nodes) > 0 {
		docs := make([]interface{}, len(nodes))
		for i, n := range nodes {
			docs[i] = n
		}
		if _, err := db.Collection("static_nodes").InsertMany(ctx, docs); err != nil {
			logger.Log.Fatal(err)
		}
	}

	if apps := inmemory.DefaultAppData(); len(apps) > 0 {
		docs := make([]interface{}, len(apps))
		for i, a := range apps {
			docs[i] = a
		}
		if _, err := db.Collection("static_apps").InsertMany(ctx, docs); err != nil {
			logger.Log.Fatal(err)
		}
	}

	mapData := inmemory.DefaultMapData()
	_, err = db.Collection("map").InsertOne(ctx, bson.M{"_id": 1, "nodes": mapData.Nodes, "connections": mapData.Connections})
	if err != nil {
		logger.Log.Fatal(err)
	}

	// create default admin account
	admin := bson.M{
		"id":          "1",
		"username":    "admin",
		"password":    "admin",
		"email":       "admin@ronaqci.eu",
		"fullname":    "Administrator",
		"affiliation": "RoNaQCI",
		"role":        "admin",
	}
	if _, err := db.Collection("auth_users").InsertOne(ctx, admin); err != nil {
		logger.Log.Fatal(err)
	}

	alertData := inmemory.DefaultAlertData()
	_, err = db.Collection("alerts_response").InsertOne(ctx, bson.M{"_id": 1, "alertlevels": alertData.AlertLevels, "alerts": alertData.Alerts})
	if err != nil {
		logger.Log.Fatal(err)
	}

	logger.Log.Info("Database populated with default data")
}
