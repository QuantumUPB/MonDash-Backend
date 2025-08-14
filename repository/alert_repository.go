package repository

import "mondash-backend/domain"

// AlertRepository defines persistence methods for alerts.
type AlertRepository interface {
	// List returns the alert configuration (levels and registered alerts).
	List() (domain.AlertInfo, error)
	Add(alert domain.Alert) error
}
