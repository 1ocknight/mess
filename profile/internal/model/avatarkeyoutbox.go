package model

import "time"

type AvatarKeyOutbox struct {
	SubjectID string
	Key       string
	CreatedAt time.Time
	DeletedAt *time.Time
}

func GetAvatarKeys(arr []*AvatarKeyOutbox) []string {
	res := make([]string, len(arr))
	for i, k := range arr {
		res[i] = k.Key
	}

	return res
}
