package domain

// App represents an application running on a node.
type App struct {
	NodeID       string `json:"nodeId"`
	Name         string `json:"name"`
	NumberOfKeys int    `json:"numberOfKeys"`
	KeySize      int    `json:"keySize"`
	Timestamp    string `json:"timestamp"`
}
