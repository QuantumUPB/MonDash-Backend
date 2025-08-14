package domain

type Alert struct {
	ID            string `json:"id"`
	Device        string `json:"device"`
	Level         string `json:"level"`
	LastActivated string `json:"lastActivated"`
	Email         string `json:"email"`
}

// AlertInfo represents the alert configuration stored by the alert service.
// It does not contain any device information.
type AlertInfo struct {
	AlertLevels []string `json:"alertLevels"`
	Alerts      []Alert  `json:"alerts"`
}

// AlertsResponse is returned by the alerts API endpoint. It combines the
// device list obtained from the device service with alert information from the
// alert service.
type AlertsResponse struct {
	Devices     []string `json:"devices"`
	AlertLevels []string `json:"alertLevels"`
	Alerts      []Alert  `json:"alerts"`
}
