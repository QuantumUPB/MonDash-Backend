package services

import "mondash-backend/domain"
import "mondash-backend/repository"

// UserService contains business logic for users.
type UserService struct {
	Repo repository.UserRepository
}

// List returns users from the repository.
func (s *UserService) List() ([]domain.User, error) {
	if s.Repo == nil {
		return nil, nil
	}
	return s.Repo.List()
}
