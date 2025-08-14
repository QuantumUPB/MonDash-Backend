package repository

import "mondash-backend/domain"

// AppRepository defines persistence methods for apps.
type AppRepository interface {
	Update(app *domain.App) error
	List() ([]domain.AppData, error)
	// Timeline returns key consumption history for all apps within the given time range.
	Timeline(start, end string) ([]domain.AppData, error)
}
