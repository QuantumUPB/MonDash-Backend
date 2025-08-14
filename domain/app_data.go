package domain

// KeyConsumptionEntry represents a single key consumption data point.
// Timestamp follows the same format used by node timestamps.
type KeyConsumptionEntry struct {
	Timestamp string `json:"timestamp"`
	Count     int    `json:"count"`
}

// AppData holds runtime information about an application.
// Additional fields may be added but the following are always present when
// returned from the /api/apps endpoint.
type AppData struct {
	Name                  string                `json:"name"`
	Certificate           string                `json:"certificate"`
	Nodes                 []string              `json:"nodes"`
	KeyConsumptionHistory []KeyConsumptionEntry `json:"keyConsumptionHistory"`
	ErrorHistory          []string              `json:"errorHistory"`
	NumberOfKeys          int                   `json:"-"`
	KeySize               int                   `json:"keySize"`
}
