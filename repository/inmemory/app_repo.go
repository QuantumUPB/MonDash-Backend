package inmemory

import (
	"errors"

	"mondash-backend/config"
	"mondash-backend/domain"
)

// AppRepo is an in-memory implementation of repository.AppRepository.
type AppRepo struct {
	data []domain.AppData
}

// DefaultAppData loads app names from the configuration file and returns them as
// AppData structures. Any parsing error results in an empty slice.
func DefaultAppData() []domain.AppData {
	cfg, err := config.LoadFromEnv()
	if err != nil {
		return []domain.AppData{}
	}
	nodesByConsumer := cfg.NodesByConsumer()
	apps := make([]domain.AppData, len(cfg.Consumers))
	for i, name := range cfg.Consumers {
		nodes := nodesByConsumer[name]
		if nodes == nil {
			nodes = []string{}
		}
		apps[i] = domain.AppData{
			Name:                  name,
			Certificate:           name + "-cert.pem",
			Nodes:                 nodes,
			KeyConsumptionHistory: []domain.KeyConsumptionEntry{},
			ErrorHistory:          []string{},
			NumberOfKeys:          0,
			KeySize:               0,
		}
	}
	return apps
}

// NewAppRepo creates a new AppRepo with data loaded from the config file.
func NewAppRepo() *AppRepo {
	return &AppRepo{data: DefaultAppData()}
}

// Update performs a basic validation and pretends to update the app.
func (r *AppRepo) Update(a *domain.App) error {
	if a == nil || a.Name == "" {
		return errors.New("invalid app")
	}
	for i := range r.data {
		if r.data[i].Name == a.Name {
			r.data[i].NumberOfKeys = a.NumberOfKeys
			r.data[i].KeySize = a.KeySize
			break
		}
	}
	return nil
}

// List returns all apps.
func (r *AppRepo) List() ([]domain.AppData, error) {
	return r.data, nil
}

// Timeline returns no data for the in-memory repository.
func (r *AppRepo) Timeline(start, end string) ([]domain.AppData, error) {
	return []domain.AppData{}, nil
}
