package inmemory

import "mondash-backend/domain"

// UserRepo is an in-memory implementation of repository.UserRepository
// backed by the in-memory AuthRepo.
// It exposes only public user information.
type UserRepo struct {
	auth *AuthRepo
}

// NewUserRepo creates a new UserRepo using the provided AuthRepo.
func NewUserRepo(auth *AuthRepo) *UserRepo {
	return &UserRepo{auth: auth}
}

// List returns all users from the AuthRepo.
func (r *UserRepo) List() ([]domain.User, error) {
	if r.auth == nil {
		return nil, nil
	}
	users := make([]domain.User, 0, len(r.auth.users))
	for _, u := range r.auth.users {
		users = append(users, u.User)
	}
	return users, nil
}
