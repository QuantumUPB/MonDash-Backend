package mongo

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"mondash-backend/domain"
	"mondash-backend/logger"
	"mondash-backend/repository"
)

// AlertRepo implements repository.AlertRepository backed by MongoDB.
type AlertRepo struct {
	coll *mongo.Collection
}

// NewAlertRepo returns a new MongoDB AlertRepo using the given database.
func NewAlertRepo(db *mongo.Database) *AlertRepo {
	return &AlertRepo{coll: db.Collection("alerts_response")}
}

// List returns the alerts response document.
func (r *AlertRepo) List() (domain.AlertInfo, error) {
	logger.Log.Debug("mongo list alerts")
	var res domain.AlertInfo
	err := r.coll.FindOne(context.Background(), bson.M{"_id": 1}).Decode(&res)
	if err == nil {
		logger.Log.Debugw("mongo list alerts result", "alerts", res)
	}
	return res, err
}

// Add pushes a new alert into the alerts array.
func (r *AlertRepo) Add(a domain.Alert) error {
	if a.ID == "" || a.Device == "" || a.Level == "" {
		return errors.New("invalid alert")
	}
	logger.Log.Debugw("mongo add alert", "id", a.ID, "device", a.Device)
	_, err := r.coll.UpdateOne(
		context.Background(),
		bson.M{"_id": 1},
		bson.M{"$push": bson.M{"alerts": a}},
		options.Update().SetUpsert(true),
	)
	return err
}

var _ repository.AlertRepository = (*AlertRepo)(nil)
