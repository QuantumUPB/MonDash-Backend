package inmemory

import (
	"sort"
	"time"

	"mondash-backend/domain"
	"mondash-backend/repository"
)

// DeviceRepo is an in-memory implementation of repository.DeviceRepository that
// aggregates devices from the NodeRepo.
type DeviceRepo struct {
	nodes   *NodeRepo
	history map[string][]domain.KeyRateEntry
}

// NewDeviceRepo creates a new DeviceRepo using the provided NodeRepo.
func NewDeviceRepo(nodeRepo *NodeRepo) *DeviceRepo {
	return &DeviceRepo{nodes: nodeRepo, history: map[string][]domain.KeyRateEntry{}}
}

// List returns all devices from all nodes.
func (r *DeviceRepo) List(silent bool) ([]domain.Device, error) {
	if r.nodes == nil {
		return []domain.Device{}, nil
	}
	nodes, err := r.nodes.List()
	if err != nil {
		return nil, err
	}
	var devices []domain.Device
	for _, n := range nodes {
		devices = append(devices, n.Devices...)
	}
	if devices == nil {
		devices = []domain.Device{}
	}
	return devices, nil
}

// KeyRateHistory returns stored key rate entries for the given device.
func (r *DeviceRepo) KeyRateHistory(id string, limit int) ([]domain.KeyRateEntry, error) {
	entries := r.history[id]
	var filtered []domain.KeyRateEntry
	for _, e := range entries {
		if e.Rate != 0 {
			filtered = append(filtered, e)
		}
	}

	sort.SliceStable(filtered, func(i, j int) bool { return filtered[i].Timestamp < filtered[j].Timestamp })

	const tolerance = 100 * time.Millisecond
	var result []domain.KeyRateEntry
	var group []domain.KeyRateEntry
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
		for _, e := range group {
			sum += e.Rate
		}
		avg := sum / len(group)
		result = append(result, domain.KeyRateEntry{Timestamp: group[0].Timestamp, Rate: avg})
		group = group[:0]
	}

	for _, rec := range filtered {
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
		result = []domain.KeyRateEntry{}
	}
	return result, nil
}

// AddKeyRate appends a key rate entry to the device's history.
func (r *DeviceRepo) AddKeyRate(id string, entry domain.KeyRateEntry) error {
	r.history[id] = append(r.history[id], entry)
	return nil
}

var _ repository.DeviceRepository = (*DeviceRepo)(nil)
