package openidconnect

type User struct {
	ID    string `json:"sub"`
	Name  string `json:"preferred_username"`
	Email string `json:"email"`
}

func (s *User) GetID() string {
	return s.Name
}

func (s *User) GetName() string {
	return s.Name
}

func (s *User) GetEmail() string {
	return s.Email
}
