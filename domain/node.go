package domain

// Node represents node runtime information sent by the agents.
type Node struct {
	Name           string  `json:"name"`
	Status         string  `json:"status"`
	StoredKeyCount int     `json:"stored_key_count"`
	CurrentKeyRate float64 `json:"current_key_rate"`
	Timestamp      string  `json:"timestamp"`
}
