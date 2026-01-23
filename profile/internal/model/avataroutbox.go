package model

import "time"

type AvatarOutbox struct {
	SubjectID string
	CreatedAt time.Time
	DeletedAt *time.Time
}

func GetOutboxIDs(arr []*AvatarOutbox) []string {
	res := make([]string, len(arr))
	for i, k := range arr {
		res[i] = k.SubjectID
	}

	return res
}
