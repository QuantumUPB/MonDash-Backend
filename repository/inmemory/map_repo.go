package inmemory

import (
	"sort"

	"mondash-backend/config"
	"mondash-backend/domain"
)

// MapRepo is an in-memory implementation of repository.MapRepository.
type MapRepo struct {
	data domain.MapData
}

// DefaultMapData loads map information from the configuration file specified by
// CONFIG_FILE and converts it to MapData. Any error results in an empty map.
func DefaultMapData() domain.MapData {
	cfg, err := config.LoadFromEnv()
	if err != nil {
		return domain.MapData{}
	}

	nodeSet := map[string]struct{}{}
	for _, name := range cfg.Names {
		nodeSet[baseName(name)] = struct{}{}
	}
	for name := range cfg.Geolocation {
		nodeSet[name] = struct{}{}
	}

	nodes := make([]domain.MapNode, 0, len(nodeSet))
	for name := range nodeSet {
		coord := cfg.Geolocation[name]
		nodes = append(nodes, domain.MapNode{
			ID:   name,
			Name: name,
			Coordinates: domain.Coordinates{
				Lat:  coord.Lat,
				Long: coord.Long,
			},
			Endpoint: true,
		})
	}
	sort.Slice(nodes, func(i, j int) bool { return nodes[i].ID < nodes[j].ID })

	var connections []domain.MapConnection
	for _, l := range cfg.Links {
		if len(l) < 2 {
			continue
		}
		from := baseName(l[0])
		to := baseName(l[1])
		connections = append(connections, domain.MapConnection{From: from, To: to, Status: "green"})
	}

	return domain.MapData{Nodes: nodes, Connections: connections}
}

// NewMapRepo creates a new MapRepo with sample data.
func NewMapRepo() *MapRepo {
	return &MapRepo{data: DefaultMapData()}
}

// Get returns the network map data.
func (r *MapRepo) Get() (domain.MapData, error) {
	return r.data, nil
}
