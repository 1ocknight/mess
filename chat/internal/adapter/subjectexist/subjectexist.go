package subjectexist

import "context"

type Service interface {
	SubjectExists(ctx context.Context, subjectID string) (bool, error)
}
