package subjectexist

import "context"

type Service interface {
	Exists(ctx context.Context, subjectID string) (bool, error)
}
