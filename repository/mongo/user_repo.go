package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"mondash-backend/domain"
	"mondash-backend/logger"
	"mondash-backend/repository"
)

// UserRepo implements repository.UserRepository backed by MongoDB.
type UserRepo struct {
	coll *mongo.Collection
}

// NewUserRepo returns a new MongoDB UserRepo using the given database.
func NewUserRepo(db *mongo.Database) *UserRepo {
	coll := db.Collection("auth_users")
	return &UserRepo{coll: coll}
}

// List returns all users from the collection.
func (r *UserRepo) List() ([]domain.User, error) {
	logger.Log.Debug("mongo list users")
	cursor, err := r.coll.Find(context.Background(), bson.D{})
	if err != nil {
		return nil, err
	}
	var users []domain.User
	if err = cursor.All(context.Background(), &users); err != nil {
		return nil, err
	}
	logger.Log.Debugw("mongo list users result", "users", users)
	return users, nil
}

var _ repository.UserRepository = (*UserRepo)(nil)
