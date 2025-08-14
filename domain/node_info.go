package domain

type Coordinates struct {
	Lat  float64 `json:"lat"`
	Long float64 `json:"long"`
}

type Connection struct {
	Device    string `json:"device"`
	OtherNode string `json:"otherNode"`
}

type Requests struct {
	TotalOverTime int            `json:"totalOverTime"`
	ByApp         map[string]int `json:"byApp"`
}

type Maintenance struct {
	From string `json:"from"`
	To   string `json:"to"`
}

type NodeInfo struct {
	ID                   string        `json:"id"`
	Name                 string        `json:"name"`
	KME                  string        `json:"kme"`
	Coordinates          Coordinates   `json:"coordinates"`
	Type                 string        `json:"type"`
	Status               string        `json:"status"`
	Connections          []Connection  `json:"connections"`
	Apps                 []string      `json:"apps"`
	Requests             Requests      `json:"requests"`
	ScheduledMaintenance []Maintenance `json:"scheduledMaintenance"`
	Devices              []Device      `json:"devices"`
	Events               []NodeEvent   `json:"events"`
}

// NodeEvent represents a status change event for a node.
type NodeEvent struct {
	Timestamp string `json:"timestamp"`
	Message   string `json:"message"`
}
