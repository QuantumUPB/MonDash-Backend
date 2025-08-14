package mongo

import (
	"context"
	"sort"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"mondash-backend/domain"
	"mondash-backend/logger"
	"mondash-backend/repository"
)

// DeviceRepo implements repository.DeviceRepository backed by MongoDB.
type DeviceRepo struct {
	coll *mongo.Collection
}

const keyRateHistoryLimit = 10

// NewDeviceRepo returns a new MongoDB DeviceRepo using the given database.
func NewDeviceRepo(db *mongo.Database) *DeviceRepo {
	return &DeviceRepo{coll: db.Collection("static_nodes")}
}

// List returns all devices from the collection.
func (r *DeviceRepo) List(silent bool) ([]domain.Device, error) {
	if !silent {
		logger.Log.Debug("mongo list devices")
	}
	cursor, err := r.coll.Find(context.Background(), bson.D{})
	if err != nil {
		return nil, err
	}
	var nodes []domain.NodeInfo
	if err = cursor.All(context.Background(), &nodes); err != nil {
		return nil, err
	}
	var devices []domain.Device
	for _, n := range nodes {
		devices = append(devices, n.Devices...)
	}
	if devices == nil {
		devices = []domain.Device{}
	}
	if !silent {
		logger.Log.Debugw("mongo list devices result", "devices", devices)
	}
	return devices, nil
}

// KeyRateHistory returns key rate history for the device ordered chronologically.
func (r *DeviceRepo) KeyRateHistory(id string, limit int) ([]domain.KeyRateEntry, error) {
	if limit <= 0 {
		limit = keyRateHistoryLimit
	}

	cursor, err := r.coll.Database().Collection("device_keyrate").Find(
		context.Background(),
		bson.M{"id": id, "rate": bson.M{"$ne": 0}},
		options.Find().SetSort(bson.M{"timestamp": -1}),
	)
	if err != nil {
		return nil, err
	}

	var recs []domain.KeyRateEntry
	if err := cursor.All(context.Background(), &recs); err != nil {
		return nil, err
	}

	sort.SliceStable(recs, func(i, j int) bool { return recs[i].Timestamp < recs[j].Timestamp })

	const tolerance = 100 * time.Millisecond
	var result []domain.KeyRateEntry
	var group []domain.KeyRateEntry
	var lastTime time.Time

	parseTS := func(ts string) (time.Time, error) {
		t, err := time.Parse(time.RFC3339Nano, ts)
		if err != nil {
			return time.Parse(time.RFC3339, ts)
		}
		return t, nil
	}

	flush := func() {
		if len(group) == 0 {
			return
		}
		sum := 0
		for _, e := range group {
			sum += e.Rate
		}
		avg := sum / len(group)
		result = append(result, domain.KeyRateEntry{Timestamp: group[0].Timestamp, Rate: avg})
		group = group[:0]
	}

	for _, rec := range recs {
		ts, err := parseTS(rec.Timestamp)
		if err != nil {
			continue
		}
		if len(group) == 0 {
			group = append(group, rec)
			lastTime = ts
			continue
		}
		if ts.Sub(lastTime) <= tolerance {
			group = append(group, rec)
		} else {
			flush()
			group = append(group, rec)
			lastTime = ts
		}
	}
	flush()

	if len(result) > limit {
		result = result[len(result)-limit:]
	}

	return result, nil
}

// AddKeyRate inserts a key rate entry for the device into the database.
func (r *DeviceRepo) AddKeyRate(id string, entry domain.KeyRateEntry) error {
	_, err := r.coll.Database().Collection("device_keyrate").InsertOne(
		context.Background(),
		bson.M{"id": id, "timestamp": entry.Timestamp, "rate": entry.Rate},
	)
	return err
}

var _ repository.DeviceRepository = (*DeviceRepo)(nil)
