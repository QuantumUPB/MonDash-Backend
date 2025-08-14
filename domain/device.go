package domain

type ConnectedTo struct {
	ID     string `json:"id"`
	NodeID string `json:"node_id"`
}

// KeyRateEntry represents a single key rate measurement for a device.
type KeyRateEntry struct {
	Timestamp string `json:"timestamp"`
	Rate      int    `json:"rate"`
}

// SelfReporting holds the runtime statistics reported by a device.
type SelfReporting struct {
	KeyRateHistory []KeyRateEntry `json:"keyrate"`
	MaxKeyRate     int            `json:"max_key_rate"`
	GenRate        int            `json:"gen_rate"`
	UsageRate      int            `json:"usage_rate"`
	Logs           []string       `json:"logs"`
}

type Device struct {
	ID            string        `json:"id"`
	Device        string        `json:"device"`
	Status        string        `json:"status"`
	NodeID        string        `json:"node_id"`
	Coordinates   Coordinates   `json:"coordinates"`
	ConnectedTo   ConnectedTo   `json:"connected_to"`
	SelfReporting SelfReporting `json:"self_reporting"`
}
