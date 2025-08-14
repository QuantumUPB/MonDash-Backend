package inmemory

import (
	"errors"

	"mondash-backend/domain"
)

// AlertRepo is an in-memory implementation of repository.AlertRepository.
type AlertRepo struct {
	data domain.AlertInfo
}

var defaultAlertData = domain.AlertInfo{
	AlertLevels: []string{"low", "medium", "high"},
	Alerts:      buildDefaultAlerts(),
}

// DefaultAlertData returns the initial alert data used by the in-memory repository.
func DefaultAlertData() domain.AlertInfo {
	return defaultAlertData
}

func buildDefaultAlerts() []domain.Alert {
	nodes := DefaultNodeData()
	var alerts []domain.Alert
	for _, n := range nodes {
		for _, d := range n.Devices {
			alerts = append(alerts, domain.Alert{
				ID:     d.ID,
				Device: d.ID,
				Level:  "high",
				Email:  "admin@ronaqci.eu",
			})
		}
	}
	return alerts
}

// NewAlertRepo creates a new AlertRepo with sample data.
func NewAlertRepo() *AlertRepo {
	return &AlertRepo{data: DefaultAlertData()}
}

// List returns all alerts.
func (r *AlertRepo) List() (domain.AlertInfo, error) {
	return r.data, nil
}

// Add appends an alert after validation.
func (r *AlertRepo) Add(a domain.Alert) error {
	if a.ID == "" || a.Device == "" || a.Level == "" {
		return errors.New("invalid alert")
	}
	r.data.Alerts = append(r.data.Alerts, a)
	return nil
}
