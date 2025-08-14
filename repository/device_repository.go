package repository

import "mondash-backend/domain"

// DeviceRepository defines persistence methods for devices.
type DeviceRepository interface {
	// List returns all devices. When silent is true, implementations should
	// suppress verbose logging.
	List(silent bool) ([]domain.Device, error)
	// KeyRateHistory returns up to `limit` key rate entries for a device.
	KeyRateHistory(deviceID string, limit int) ([]domain.KeyRateEntry, error)
	// AddKeyRate stores a new key rate entry for the given device.
	AddKeyRate(deviceID string, entry domain.KeyRateEntry) error
}
