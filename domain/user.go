package domain

type User struct {
	ID          string `json:"id"`
	Email       string `json:"email"`
	FullName    string `json:"fullName"`
	Affiliation string `json:"affiliation"`
	Role        string `json:"role"`
}
