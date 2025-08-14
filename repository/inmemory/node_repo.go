package inmemory

import (
	"errors"
	"sort"

	"mondash-backend/config"
	"mondash-backend/domain"
)

// NodeRepo is an in-memory implementation of repository.NodeRepository.
type NodeRepo struct {
	data []domain.NodeInfo
}

// DefaultNodeData loads node data from the configuration file defined by
// CONFIG_FILE. If the file cannot be read, an empty slice is returned.
func DefaultNodeData() []domain.NodeInfo {
	cfg, err := config.LoadFromEnv()
	if err != nil {
		return []domain.NodeInfo{}
	}

	groups := map[string][]string{}
	for _, name := range cfg.Names {
		base := baseName(name)
		groups[base] = append(groups[base], name)
	}

	consumersByNode := cfg.ConsumersByNode()

	var (
		nodes     []domain.NodeInfo
		deviceMap = make(map[string]*domain.Device)
	)
	for name, devs := range groups {
		coord := cfg.Geolocation[name]
		apps := consumersByNode[name]
		if apps == nil {
			apps = []string{}
		}
		nodeType := "terminal"
		if len(devs) > 1 {
			nodeType = "trusted node"
		}
		node := domain.NodeInfo{
			ID:     name,
			Name:   name,
			Status: "active",
			KME:    name + "-kme",
			Type:   nodeType,
			Apps:   apps,
			Events: []domain.NodeEvent{},
			Coordinates: domain.Coordinates{
				Lat:  coord.Lat,
				Long: coord.Long,
			},
		}
		for _, d := range devs {
			device := domain.Device{
				ID:          d,
				Device:      d,
				Status:      "online",
				NodeID:      name,
				Coordinates: node.Coordinates,
			}
			node.Devices = append(node.Devices, device)
		}
		for _, l := range cfg.Links {
			if len(l) < 2 {
				continue
			}
			from := l[0]
			to := l[1]
			if baseName(from) == name {
				node.Connections = append(node.Connections, domain.Connection{
					Device:    from,
					OtherNode: baseName(to),
				})
			}
			if baseName(to) == name {
				node.Connections = append(node.Connections, domain.Connection{
					Device:    to,
					OtherNode: baseName(from),
				})
			}
		}
		nodes = append(nodes, node)
	}
	// build lookup map from device ID to pointer for setting connections
	for i := range nodes {
		for j := range nodes[i].Devices {
			deviceMap[nodes[i].Devices[j].ID] = &nodes[i].Devices[j]
		}
	}
	// assign ConnectedTo information based on link configuration
	for _, l := range cfg.Links {
		if len(l) < 2 {
			continue
		}
		from, to := l[0], l[1]
		if d := deviceMap[from]; d != nil {
			d.ConnectedTo = domain.ConnectedTo{ID: to, NodeID: baseName(to)}
		}
		if d := deviceMap[to]; d != nil {
			d.ConnectedTo = domain.ConnectedTo{ID: from, NodeID: baseName(from)}
		}
	}
	sort.Slice(nodes, func(i, j int) bool { return nodes[i].ID < nodes[j].ID })
	return nodes
}

func baseName(s string) string {
	if len(s) == 0 {
		return s
	}
	r := s[len(s)-1]
	if r >= 'A' && r <= 'Z' {
		return s[:len(s)-1]
	}
	return s
}

// NewNodeRepo creates a new NodeRepo with data loaded from the config file.
func NewNodeRepo() *NodeRepo {
	return &NodeRepo{data: DefaultNodeData()}
}

// Update performs validation and pretends to update a node.
func (r *NodeRepo) Update(nodes []domain.Node) error {
	if len(nodes) == 0 {
		return errors.New("invalid nodes")
	}
	for _, n := range nodes {
		if n.Name == "" {
			return errors.New("invalid node")
		}
		if n.Timestamp == "" {
			return errors.New("missing timestamp")
		}
		for i := range r.data {
			if r.data[i].Name != n.Name {
				continue
			}
			prev := r.data[i].Status
			if prev != n.Status {
				var msg string
				switch {
				case prev != "down" && n.Status == "down":
					msg = "node went down"
				case prev == "down" && n.Status != "down":
					msg = "node went up"
				}
				if msg != "" {
					r.data[i].Events = append(r.data[i].Events, domain.NodeEvent{Timestamp: n.Timestamp, Message: msg})
				}
				r.data[i].Status = n.Status
			}
			break
		}
	}
	return nil
}

// List returns all nodes.
func (r *NodeRepo) List() ([]domain.NodeInfo, error) {
	return r.data, nil
}
