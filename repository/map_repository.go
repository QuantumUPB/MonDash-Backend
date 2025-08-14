package repository

import "mondash-backend/domain"

// MapRepository defines persistence methods for network map data.
type MapRepository interface {
	Get() (domain.MapData, error)
}
