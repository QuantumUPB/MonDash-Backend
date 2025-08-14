package api

import (
	"encoding/json"
	"net/http"
	"os"
	"strconv"

	"mondash-backend/domain"
	"mondash-backend/services"
)

// AppsHandler returns app information via the service.
func AppsHandler(s *services.AppService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := s.List()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(data)
	}
}

// AppsTimelineHandler returns app keyrate information within a time range via the service.
func AppsTimelineHandler(s *services.AppService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := r.URL.Query().Get("startTimestamp")
		end := r.URL.Query().Get("endTimestamp")
		data, err := s.Timeline(start, end)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(data)
	}
}

// AlertsHandler returns device information via the service.
// This mirrors the behaviour of the /api/devices endpoint.
// AlertsHandler returns the available devices along with alert configuration.
// The device list is fetched from the device service while alert levels and
// registered alerts are provided by the alert service.
func AlertsHandler(d *services.DeviceService, a *services.AlertService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		devices, err := d.List()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		alertData, err := a.List()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var names []string
		for _, dev := range devices {
			if dev.Device != "" {
				names = append(names, dev.Device)
			} else if dev.ID != "" {
				names = append(names, dev.ID)
			}
		}

		json.NewEncoder(w).Encode(domain.AlertsResponse{
			Devices:     names,
			AlertLevels: alertData.AlertLevels,
			Alerts:      alertData.Alerts,
		})
	}
}

// RegisterAlertHandler accepts an alert registration and stores it.
func RegisterAlertHandler(s *services.AlertService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Device string `json:"device"`
			Level  string `json:"level"`
			Email  string `json:"email"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if err := s.Register(domain.Alert{Device: req.Device, Level: req.Level, Email: req.Email}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("ok"))
	}
}

// ActiveAlertsHandler returns the list of currently active alerts.
func ActiveAlertsHandler(s *services.AlertService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := s.ActiveAlerts()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(struct {
			Alerts []domain.Alert `json:"alerts"`
		}{Alerts: data})
	}
}

// NodesHandler returns node information via the service.
func NodesHandler(s *services.NodeService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := s.List()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(data)
	}
}

// MapHandler returns network map information via the service.
func MapHandler(s *services.MapService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := s.Get()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(data)
	}
}

// DevicesHandler returns device information via the service.
func DevicesHandler(s *services.DeviceService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		limit := 0
		if v := r.URL.Query().Get("numEntries"); v != "" {
			if n, err := strconv.Atoi(v); err == nil {
				limit = n
			}
		}
		data, err := s.ListWithHistory(limit)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(data)
	}
}

// UsersHandler returns user information via the service.
func UsersHandler(s *services.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := s.List()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(data)
	}
}

// LoginHandler accepts credentials and returns a token via the service.
func LoginHandler(s *services.AuthService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		token, err := s.Login(req.Username, req.Password)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		if token == "" {
			token = os.Getenv("AUTH_TOKEN")
			if token == "" {
				token = "abc"
			}
		}
		http.SetCookie(w, &http.Cookie{
			Name:     "auth_token",
			Value:    token,
			Path:     "/",
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteNoneMode,
		})
		json.NewEncoder(w).Encode(struct {
			Token string `json:"token"`
		}{Token: token})
	}
}

// RegisterHandler accepts registration info and delegates to the service.
func RegisterHandler(s *services.AuthService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Username string `json:"username"`
			Email    string `json:"email"`
			Password string `json:"password"`
			Role     string `json:"role"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if err := s.Register(req.Username, req.Email, req.Password, req.Role); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("ok"))
	}
}
