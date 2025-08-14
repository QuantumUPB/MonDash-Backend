package services

import "mondash-backend/repository"
import "mondash-backend/domain"

// MapService contains business logic for map data.
type MapService struct {
	Repo repository.MapRepository
}

// Get returns map data from the repository.
func (s *MapService) Get() (domain.MapData, error) {
	if s.Repo == nil {
		return domain.MapData{}, nil
	}
	return s.Repo.Get()
}
