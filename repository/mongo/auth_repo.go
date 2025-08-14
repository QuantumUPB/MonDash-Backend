package mongo

import (
	"context"
	"errors"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"mondash-backend/logger"
	"mondash-backend/repository"
)

// AuthRepo implements repository.AuthRepository backed by MongoDB.
type AuthRepo struct {
	coll *mongo.Collection
}

// NewAuthRepo returns a new MongoDB AuthRepo using the given database.
func NewAuthRepo(db *mongo.Database) *AuthRepo {
	coll := db.Collection("auth_users")
	repo := &AuthRepo{coll: coll}

	// ensure a default admin account exists
	ctx := context.Background()
	count, err := coll.CountDocuments(ctx, bson.D{})
	if err == nil && count == 0 {
		admin := bson.M{
			"id":          "1",
			"username":    "admin",
			"password":    "admin",
			"email":       "admin@ronaqci.eu",
			"fullname":    "Administrator",
			"affiliation": "RoNaQCI",
			"role":        "admin",
		}
		if _, err := coll.InsertOne(ctx, admin); err != nil {
			logger.Log.Errorw("failed to insert default admin account", "error", err)
		}
	}

	return repo
}

// Login checks credentials and returns a dummy token if valid.
func (r *AuthRepo) Login(username, password string) (string, error) {
	if username == "" || password == "" {
		return "", errors.New("missing credentials")
	}
	logger.Log.Debugw("mongo login", "username", username)
	var doc bson.M
	err := r.coll.FindOne(context.Background(), bson.M{"username": username, "password": password}).Decode(&doc)
	if err != nil {
		return "", errors.New("invalid credentials")
	}
	if doc != nil {
		sanitized := bson.M{"username": doc["username"]}
		logger.Log.Debugw("mongo login result", "doc", sanitized)
	}
	token := os.Getenv("AUTH_TOKEN")
	if token == "" {
		token = "abc"
	}
	return token, nil
}

// Register inserts a new user document.
func (r *AuthRepo) Register(username, email, password, role string) error {
	if username == "" || email == "" || password == "" || role == "" {
		return errors.New("invalid registration")
	}
	logger.Log.Debugw("mongo register user", "username", username)
	doc := bson.M{
		"username":    username,
		"password":    password,
		"email":       email,
		"fullname":    username,
		"affiliation": "",
		"role":        role,
	}
	_, err := r.coll.InsertOne(context.Background(), doc)
	return err
}

var _ repository.AuthRepository = (*AuthRepo)(nil)
