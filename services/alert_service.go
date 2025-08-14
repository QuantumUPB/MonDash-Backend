package services

import (
	"context"
	"fmt"
	"net/smtp"
	"os"
	"strings"
	"time"

	"mondash-backend/domain"
	"mondash-backend/logger"
	"mondash-backend/repository"
)

// AlertService contains business logic for alerts.
type AlertService struct {
	Repo       repository.AlertRepository
	DeviceRepo repository.DeviceRepository

	registered []domain.Alert

	emailEnabled bool
	smtpHost     string
	smtpPort     string
	smtpUser     string
	smtpPass     string
	smtpFrom     string
}

// InitFromEnv loads email settings from environment variables.
func (s *AlertService) InitFromEnv() {
	s.emailEnabled = strings.ToLower(os.Getenv("EMAIL_ON_ALERT")) == "true"
	s.smtpHost = os.Getenv("SMTP_HOST")
	s.smtpPort = os.Getenv("SMTP_PORT")
	s.smtpUser = os.Getenv("SMTP_USERNAME")
	s.smtpPass = os.Getenv("SMTP_PASSWORD")
	s.smtpFrom = os.Getenv("SMTP_FROM")
}

// List returns alerts from the repository.
func (s *AlertService) List() (domain.AlertInfo, error) {
	if s.Repo == nil {
		return domain.AlertInfo{}, nil
	}
	return s.Repo.List()
}

// Load fetches the registered alerts from the repository.
func (s *AlertService) Load() error {
	if s.Repo == nil {
		return nil
	}
	res, err := s.Repo.List()
	if err != nil {
		return err
	}
	s.registered = res.Alerts
	return nil
}

// Register stores a new alert.
func (s *AlertService) Register(a domain.Alert) error {
	if s.Repo == nil {
		return nil
	}
	if a.ID == "" {
		a.ID = fmt.Sprintf("alert-%d", time.Now().UnixNano())
	}
	if err := s.Repo.Add(a); err != nil {
		return err
	}
	s.registered = append(s.registered, a)
	return nil
}

// ActiveAlerts returns all alerts whose device status is offline.
func (s *AlertService) ActiveAlerts() ([]domain.Alert, error) {
	if s.DeviceRepo == nil {
		return nil, nil
	}
	devices, err := s.DeviceRepo.List(false)
	if err != nil {
		return nil, err
	}
	status := map[string]string{}
	for _, d := range devices {
		status[d.ID] = d.Status
	}
	var actives []domain.Alert
	for _, a := range s.registered {
		if st, ok := status[a.Device]; ok && (st == "down" || st == "offline") {
			actives = append(actives, a)
		}
	}
	return actives, nil
}

// StartMonitoring periodically scans devices and logs down ones.
func (s *AlertService) StartMonitoring(ctx context.Context, interval time.Duration) {
	if s.DeviceRepo == nil {
		return
	}
	ticker := time.NewTicker(interval)
	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				s.scan()
			}
		}
	}()
}

func (s *AlertService) scan() {
	devices, err := s.DeviceRepo.List(true)
	if err != nil {
		return
	}
	status := make(map[string]string)
	for _, d := range devices {
		status[d.ID] = d.Status
	}

	for i, a := range s.registered {
		st, ok := status[a.Device]
		if !ok {
			continue
		}
		if st == "down" || st == "offline" {
			if a.LastActivated == "" {
				logger.Log.Infof("device %s is down", a.Device)
				_ = s.sendEmail(a.Email, "Device down", fmt.Sprintf("device %s is down", a.Device))
				a.LastActivated = time.Now().Format(time.RFC3339)
				s.registered[i] = a
			}
		} else {
			if a.LastActivated != "" {
				a.LastActivated = ""
				s.registered[i] = a
			}
		}
	}
}

func (s *AlertService) sendEmail(to, subject, body string) error {
	if !s.emailEnabled {
		return nil
	}
	if s.smtpHost == "" || s.smtpPort == "" || s.smtpFrom == "" {
		return fmt.Errorf("smtp not configured")
	}
	addr := s.smtpHost + ":" + s.smtpPort
	msg := []byte("To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" + body + "\r\n")
	var auth smtp.Auth
	if s.smtpUser != "" || s.smtpPass != "" {
		auth = smtp.PlainAuth("", s.smtpUser, s.smtpPass, s.smtpHost)
	}
	return smtp.SendMail(addr, auth, s.smtpFrom, []string{to}, msg)
}
