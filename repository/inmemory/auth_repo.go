package inmemory

import (
	"errors"
	"os"
	"strconv"

	"mondash-backend/domain"
)

// AuthUser represents a stored authentication user.
type AuthUser struct {
	domain.User
	Username string
	Password string
}

// AuthRepo is an in-memory implementation of repository.AuthRepository.
// It stores users in a map keyed by username.
type AuthRepo struct {
	users map[string]AuthUser
}

// NewAuthRepo creates a new AuthRepo seeded with a default admin user.
func NewAuthRepo() *AuthRepo {
	admin := AuthUser{
		User: domain.User{
			ID:          "1",
			Email:       "admin@ronaqci.eu",
			FullName:    "Administrator",
			Affiliation: "RoNaQCI",
			Role:        "admin",
		},
		Username: "admin",
		Password: "admin",
	}
	return &AuthRepo{users: map[string]AuthUser{"admin": admin}}
}

// Login returns a dummy token if credentials are non-empty.
func (r *AuthRepo) Login(username, password string) (string, error) {
	user, ok := r.users[username]
	if !ok || user.Password != password {
		return "", errors.New("invalid credentials")
	}
	token := os.Getenv("AUTH_TOKEN")
	if token == "" {
		token = "abc"
	}
	return token, nil
}

// Register performs basic validation.
func (r *AuthRepo) Register(username, email, password, role string) error {
	if username == "" || email == "" || password == "" || role == "" {
		return errors.New("invalid registration")
	}
	id := strconv.Itoa(len(r.users) + 1)
	r.users[username] = AuthUser{
		User: domain.User{
			ID:          id,
			Email:       email,
			FullName:    username,
			Affiliation: "",
			Role:        role,
		},
		Username: username,
		Password: password,
	}
	return nil
}
