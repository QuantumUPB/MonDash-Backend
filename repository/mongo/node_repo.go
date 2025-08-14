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

// NodeRepo implements repository.NodeRepository backed by MongoDB. Static
// information is stored in the `static_nodes` collection while dynamic updates
// are appended to the `node_history` collection.
type NodeRepo struct {
	staticColl  *mongo.Collection
	dynamicColl *mongo.Collection
}

// NewNodeRepo returns a new MongoDB NodeRepo using the given database.
func NewNodeRepo(db *mongo.Database) *NodeRepo {
	return &NodeRepo{
		staticColl:  db.Collection("static_nodes"),
		dynamicColl: db.Collection("node_history"),
	}
}

// Update upserts a node's name by ID.
func (r *NodeRepo) Update(nodes []domain.Node) error {
	if len(nodes) == 0 {
		return errors.New("invalid nodes")
	}
	for _, n := range nodes {
		logger.Log.Debugw("mongo store node update", "name", n.Name)
		if n.Name == "" {
			return errors.New("invalid node")
		}
		if n.Timestamp == "" {
			return errors.New("missing timestamp")
		}
		var last domain.Node
		err := r.dynamicColl.FindOne(
			context.Background(),
			bson.M{"name": n.Name},
			options.FindOne().SetSort(bson.M{"timestamp": -1}),
		).Decode(&last)
		if err == nil {
			var msg string
			switch {
			case last.Status != "down" && n.Status == "down":
				msg = "node went down"
			case last.Status == "down" && n.Status != "down":
				msg = "node went up"
			}
			if msg != "" {
				event := domain.NodeEvent{Timestamp: n.Timestamp, Message: msg}
				_, _ = r.staticColl.UpdateOne(
					context.Background(),
					bson.M{"name": n.Name},
					bson.M{"$push": bson.M{"events": event}},
				)
			}
		}
		_, err = r.dynamicColl.InsertOne(context.Background(), n)
		if err != nil {
			return err
		}
	}
	return nil
}

// List returns all node information from the collection.
func (r *NodeRepo) List() ([]domain.NodeInfo, error) {
	logger.Log.Debug("mongo list nodes")
	cursor, err := r.staticColl.Find(context.Background(), bson.D{})
	if err != nil {
		return nil, err
	}
	var nodes []domain.NodeInfo
	if err = cursor.All(context.Background(), &nodes); err != nil {
		return nil, err
	}
	for i := range nodes {
		var update domain.Node
		err := r.dynamicColl.FindOne(
			context.Background(),
			bson.M{"name": nodes[i].Name},
			options.FindOne().SetSort(bson.M{"timestamp": -1}),
		).Decode(&update)
		if err == nil {
			nodes[i].Status = update.Status
		}
	}
	logger.Log.Debugw("mongo list nodes result", "nodes", nodes)
	return nodes, nil
}

var _ repository.NodeRepository = (*NodeRepo)(nil)
