package repository

import "mondash-backend/domain"

// UserRepository defines persistence methods for users.
type UserRepository interface {
	List() ([]domain.User, error)
}
