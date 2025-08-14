package mongo

import (
	"context"
	"errors"
	"sort"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"mondash-backend/config"
	"mondash-backend/domain"
	"mondash-backend/logger"
	"mondash-backend/repository"
)

// AppRepo implements repository.AppRepository backed by MongoDB.
type AppRepo struct {
	staticColl  *mongo.Collection
	dynamicColl *mongo.Collection
}

func (r *AppRepo) latest(name string) (*domain.App, error) {
	var rec domain.App
	err := r.dynamicColl.FindOne(
		context.Background(),
		bson.M{"name": name},
		options.FindOne().SetSort(bson.M{"timestamp": -1}),
	).Decode(&rec)
	if err != nil {
		return nil, err
	}
	return &rec, nil
}

const historyLimit = 10

// NewAppRepo returns a new MongoDB AppRepo using the given database.
func NewAppRepo(db *mongo.Database) *AppRepo {
	return &AppRepo{
		staticColl:  db.Collection("static_apps"),
		dynamicColl: db.Collection("key_consumption"),
	}
}

// historyFor returns the last `limit` key consumption entries for the given app
// name. Results are ordered chronologically with the earliest entry first.
func (r *AppRepo) historyFor(name string, limit int) ([]domain.KeyConsumptionEntry, error) {
	if limit <= 0 {
		limit = historyLimit
	}
	cursor, err := r.dynamicColl.Find(
		context.Background(),
		bson.M{"name": name},
		options.Find().SetSort(bson.M{"timestamp": -1}),
	)
	if err != nil {
		return nil, err
	}
	var recs []domain.App
	if err := cursor.All(context.Background(), &recs); err != nil {
		return nil, err
	}

	sort.SliceStable(recs, func(i, j int) bool { return recs[i].Timestamp < recs[j].Timestamp })

	const tolerance = 100 * time.Millisecond
	var result []domain.KeyConsumptionEntry
	var group []domain.App
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
		for _, g := range group {
			if g.NumberOfKeys > 0 {
				sum += g.NumberOfKeys
			} else {
				sum++
			}
		}
		result = append(result, domain.KeyConsumptionEntry{Timestamp: group[0].Timestamp, Count: sum})
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

	if limit > 0 && len(result) > limit {
		result = result[len(result)-limit:]
	}
	if result == nil {
		result = []domain.KeyConsumptionEntry{}
	}
	return result, nil
}

// Update upserts an app's basic information.
func (r *AppRepo) Update(a *domain.App) error {
	if a == nil || a.Name == "" || a.Timestamp == "" {
		return errors.New("invalid app")
	}
	logger.Log.Debugw("mongo store app update", "name", a.Name)
	_, err := r.dynamicColl.InsertOne(context.Background(), a)
	return err
}

// List returns all app data from the collection.
func (r *AppRepo) List() ([]domain.AppData, error) {
	logger.Log.Debug("mongo list apps")
	cursor, err := r.staticColl.Find(context.Background(), bson.D{})
	if err != nil {
		return nil, err
	}
	var apps []domain.AppData
	if err = cursor.All(context.Background(), &apps); err != nil {
		return nil, err
	}

	cfg, _ := config.LoadFromEnv()
	nodesByConsumer := cfg.NodesByConsumer()

	for i := range apps {
		history, err := r.historyFor(apps[i].Name, historyLimit)
		if err == nil {
			apps[i].KeyConsumptionHistory = history
		}
		if apps[i].KeyConsumptionHistory == nil {
			apps[i].KeyConsumptionHistory = []domain.KeyConsumptionEntry{}
		}
		if apps[i].ErrorHistory == nil {
			apps[i].ErrorHistory = []string{}
		}
		if nodes := nodesByConsumer[apps[i].Name]; len(nodes) > 0 {
			apps[i].Nodes = nodes
		} else if apps[i].Nodes == nil {
			apps[i].Nodes = []string{}
		}
		if last, err := r.latest(apps[i].Name); err == nil && last != nil {
			apps[i].KeySize = last.KeySize
		}
	}
	logger.Log.Debugw("mongo list apps result", "apps", apps)
	return apps, nil
}

// Timeline returns key consumption history for all apps within the given range.
func (r *AppRepo) Timeline(start, end string) ([]domain.AppData, error) {
	filter := bson.M{}
	ts := bson.M{}
	if start != "" {
		ts["$gte"] = start
	}
	if end != "" {
		ts["$lte"] = end
	}
	if len(ts) > 0 {
		filter["timestamp"] = ts
	}
	cursor, err := r.dynamicColl.Find(
		context.Background(),
		filter,
		options.Find().SetSort(bson.M{"timestamp": 1}),
	)
	if err != nil {
		return nil, err
	}
	var recs []domain.App
	if err := cursor.All(context.Background(), &recs); err != nil {
		return nil, err
	}

	byName := make(map[string][]domain.App)
	for _, rec := range recs {
		byName[rec.Name] = append(byName[rec.Name], rec)
	}

	const tolerance = 100 * time.Millisecond
	parseTS := func(ts string) (time.Time, error) {
		t, err := time.Parse(time.RFC3339Nano, ts)
		if err != nil {
			return time.Parse(time.RFC3339, ts)
		}
		return t, nil
	}

	result := make([]domain.AppData, 0, len(byName))
	for name, entries := range byName {
		sort.Slice(entries, func(i, j int) bool { return entries[i].Timestamp < entries[j].Timestamp })

		var group []domain.App
		var lastTime time.Time
		var history []domain.KeyConsumptionEntry

		flush := func() {
			if len(group) == 0 {
				return
			}
			sum := 0
			for _, g := range group {
				if g.NumberOfKeys > 0 {
					sum += g.NumberOfKeys
				} else {
					sum++
				}
			}
			history = append(history, domain.KeyConsumptionEntry{Timestamp: group[0].Timestamp, Count: sum})
			group = group[:0]
		}

		for _, rec := range entries {
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

		if history == nil {
			history = []domain.KeyConsumptionEntry{}
		}
		result = append(result, domain.AppData{
			Name:                  name,
			Nodes:                 []string{},
			KeyConsumptionHistory: history,
			ErrorHistory:          []string{},
			KeySize:               entries[len(entries)-1].KeySize,
		})
	}

	return result, nil
}

var _ repository.AppRepository = (*AppRepo)(nil)
