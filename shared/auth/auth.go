package auth

import "github.com/1ocknight/mess/shared/model"

type Service interface {
	//gSubjectExists(id string) (bool, error)
	Verify(src string) (model.Subject, error)
}

type DeleteSubjectEvent interface {
	GetSubjectID() string
}
