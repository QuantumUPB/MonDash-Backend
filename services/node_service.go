package services

import (
	"time"

	"mondash-backend/domain"
	"mondash-backend/repository"
)

// NodeService contains business logic for nodes.
type NodeService struct {
	Repo       repository.NodeRepository
	DeviceRepo repository.DeviceRepository
}

// Update updates a node using the repository.
func (s *NodeService) Update(nodes []domain.Node) error {
	now := time.Now().Format(time.RFC3339)
	for i := range nodes {
		if nodes[i].Timestamp == "" {
			nodes[i].Timestamp = now
		}
		if s.DeviceRepo != nil {
			entry := domain.KeyRateEntry{Timestamp: nodes[i].Timestamp, Rate: int(nodes[i].CurrentKeyRate)}
			_ = s.DeviceRepo.AddKeyRate(nodes[i].Name, entry)
		}
	}
	if s.Repo == nil {
		return nil
	}
	return s.Repo.Update(nodes)
}

// List returns nodes from the repository.
func (s *NodeService) List() ([]domain.NodeInfo, error) {
	if s.Repo == nil {
		return nil, nil
	}
	nodes, err := s.Repo.List()
	if err != nil {
		return nil, err
	}
	for i := range nodes {
		if nodes[i].Status == "up" {
			nodes[i].Status = "active"
		}
		if nodes[i].Type == "trusted node" {
			seen := make(map[string]struct{})
			apps := nodes[i].Apps[:0]
			for _, app := range nodes[i].Apps {
				if _, ok := seen[app]; ok {
					continue
				}
				seen[app] = struct{}{}
				apps = append(apps, app)
			}
			nodes[i].Apps = apps
		}
	}
	return nodes, nil
}
