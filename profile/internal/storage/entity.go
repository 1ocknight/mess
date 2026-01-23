package storage

import (
	"time"

	"github.com/TATAROmangol/mess/profile/internal/model"
)

type ProfileEntity struct {
	SubjectID string     `db:"subject_id"`
	Alias     string     `db:"alias"`
	Version   int        `db:"version"`
	UpdatedAt time.Time  `db:"updated_at"`
	CreatedAt time.Time  `db:"created_at"`
	DeletedAt *time.Time `db:"deleted_at"`
}

func (p *ProfileEntity) ToModel() *model.Profile {
	return &model.Profile{
		SubjectID: p.SubjectID,
		Alias:     p.Alias,
		Version:   p.Version,
		UpdatedAt: p.UpdatedAt,
		CreatedAt: p.CreatedAt,
		DeletedAt: p.DeletedAt,
	}
}

func (p *ProfileEntity) Key() *string {
	return &p.SubjectID
}

func ProfileEntitiesToModels(entities []*ProfileEntity) []*model.Profile {
	models := make([]*model.Profile, 0, len(entities))
	for _, entity := range entities {
		models = append(models, entity.ToModel())
	}
	return models
}

type AvatarOutboxEntity struct {
	SubjectID string     `db:"subject_id"`
	CreatedAt time.Time  `db:"created_at"`
	DeletedAt *time.Time `db:"deleted_at"`
}

func (p *AvatarOutboxEntity) ToModel() *model.AvatarOutbox {
	return &model.AvatarOutbox{
		SubjectID: p.SubjectID,
		DeletedAt: p.DeletedAt,
		CreatedAt: p.CreatedAt,
	}
}

func AvatarOutboxEntitiesToModels(entities []*AvatarOutboxEntity) []*model.AvatarOutbox {
	models := make([]*model.AvatarOutbox, 0, len(entities))
	for _, entity := range entities {
		models = append(models, entity.ToModel())
	}
	return models
}
