package model

type Subject interface {
	GetSubjectId() string
	GetName() string
	GetEmail() string
}
