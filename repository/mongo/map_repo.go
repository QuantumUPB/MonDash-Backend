package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"mondash-backend/domain"
	"mondash-backend/logger"
	"mondash-backend/repository"
)

// MapRepo implements repository.MapRepository backed by MongoDB.
type MapRepo struct {
	coll *mongo.Collection
}

// NewMapRepo returns a new MongoDB MapRepo using the given database.
func NewMapRepo(db *mongo.Database) *MapRepo {
	return &MapRepo{coll: db.Collection("map")}
}

// Get returns the network map data from the collection.
func (r *MapRepo) Get() (domain.MapData, error) {
	logger.Log.Debug("mongo get map data")
	var data domain.MapData
	err := r.coll.FindOne(context.Background(), bson.D{}).Decode(&data)
	if err == nil {
		logger.Log.Debugw("mongo get map data result", "data", data)
	}
	return data, err
}

var _ repository.MapRepository = (*MapRepo)(nil)
