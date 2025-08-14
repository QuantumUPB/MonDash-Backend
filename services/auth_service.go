package services

import "mondash-backend/repository"

// AuthService contains authentication business logic.
type AuthService struct {
	Repo repository.AuthRepository
}

// Login delegates to the repository.
func (s *AuthService) Login(username, password string) (string, error) {
	if s.Repo == nil {
		return "", nil
	}
	return s.Repo.Login(username, password)
}

// Register delegates to the repository.
func (s *AuthService) Register(username, email, password, role string) error {
	if s.Repo == nil {
		return nil
	}
	return s.Repo.Register(username, email, password, role)
}
