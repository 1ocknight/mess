package model

type User struct {
	ID    string
	Name  string
	Email string
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
