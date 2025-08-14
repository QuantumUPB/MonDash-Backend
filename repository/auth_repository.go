package repository

// AuthRepository defines authentication persistence methods.
type AuthRepository interface {
	Login(username, password string) (string, error)
	Register(username, email, password, role string) error
}
