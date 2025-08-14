package mongo

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"mondash-backend/logger"
)

// Connect establishes a connection to MongoDB and returns the database.
func Connect(uri, dbName string) (*mongo.Database, error) {
	logger.Log.Infow("connecting to MongoDB", "uri", uri)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}
	return client.Database(dbName), nil
}
