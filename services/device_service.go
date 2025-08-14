package services

import "mondash-backend/domain"
import "mondash-backend/repository"

// DeviceService contains business logic for devices.
type DeviceService struct {
	Repo repository.DeviceRepository
}

// List returns devices from the repository.
func (s *DeviceService) List() ([]domain.Device, error) {
	if s.Repo == nil {
		return []domain.Device{}, nil
	}
	devices, err := s.Repo.List(false)
	if devices == nil && err == nil {
		devices = []domain.Device{}
	}
	return devices, err
}

// ListWithHistory returns devices and augments each with key rate history.
func (s *DeviceService) ListWithHistory(limit int) ([]domain.Device, error) {
	devices, err := s.List()
	if err != nil || s.Repo == nil {
		return devices, err
	}
	for i := range devices {
		history, errHist := s.Repo.KeyRateHistory(devices[i].ID, limit)
		if errHist != nil {
			devices[i].SelfReporting.KeyRateHistory = []domain.KeyRateEntry{}
			continue
		}
		if history == nil {
			history = []domain.KeyRateEntry{}
		}
		devices[i].SelfReporting.KeyRateHistory = history
		if len(history) > 0 {
			max := history[0].Rate
			for _, entry := range history {
				if entry.Rate > max {
					max = entry.Rate
				}
			}
			devices[i].SelfReporting.MaxKeyRate = max
			devices[i].SelfReporting.GenRate = history[len(history)-1].Rate
		}
	}
	return devices, nil
}
