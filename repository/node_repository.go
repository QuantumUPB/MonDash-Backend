package repository

import "mondash-backend/domain"

// NodeRepository defines persistence methods for nodes.
type NodeRepository interface {
	Update(nodes []domain.Node) error
	List() ([]domain.NodeInfo, error)
}
