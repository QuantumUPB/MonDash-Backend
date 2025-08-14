package services

import (
	"time"

	"mondash-backend/domain"
	"mondash-backend/repository"
)

// AppService contains business logic for apps.
type AppService struct {
	Repo repository.AppRepository
}

// Update updates an app using the repository.
func (s *AppService) Update(a *domain.App) error {
	if s.Repo == nil {
		return nil
	}
	if a.Timestamp == "" {
		a.Timestamp = time.Now().Format(time.RFC3339)
	}
	return s.Repo.Update(a)
}

// List returns apps from the repository.
func (s *AppService) List() ([]domain.AppData, error) {
	if s.Repo == nil {
		return nil, nil
	}
	return s.Repo.List()
}

// Timeline returns key consumption history for all apps within a time range.
func (s *AppService) Timeline(start, end string) ([]domain.AppData, error) {
	if s.Repo == nil {
		return nil, nil
	}
	return s.Repo.Timeline(start, end)
}
