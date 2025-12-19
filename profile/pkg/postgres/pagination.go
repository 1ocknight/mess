package postgres

import (
	"encoding/base64"
	"encoding/json"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

const (
	AscSortLabel  = "ASC"
	DescSortLabel = "DESC"
)

type Sort struct {
	Field string `json:"field"`
	Asc   bool   `json:"asc"`
}

type Last struct {
	Field string `json:"field"`
	Key   any    `json:"key"`
}

type Pagination struct {
	Sort *Sort `json:"sort,omitempty"`
	Last *Last `json:"last,omitempty"`
	Size int   `json:"size"`
}

func NewPagination(size int, sort *Sort, last *Last) *Pagination {
	return &Pagination{
		Sort: sort,
		Last: last,
		Size: size,
	}
}

func GetPaginationToken(p *Pagination) string {
	if p == nil {
		return ""
	}

	data, err := json.Marshal(p)
	if err != nil {
		return ""
	}

	return base64.RawURLEncoding.EncodeToString(data)
}

func ParsePaginationToken(token string) (*Pagination, error) {
	if token == "" {
		return nil, nil
	}

	raw, err := base64.RawURLEncoding.DecodeString(token)
	if err != nil {
		return nil, fmt.Errorf("decode pagination token: %w", err)
	}

	var p Pagination
	if err := json.Unmarshal(raw, &p); err != nil {
		return nil, fmt.Errorf("unmarshal pagination token: %w", err)
	}

	return &p, nil
}

func MakeQueryWithPagination[T any](db *sqlx.DB, b sq.SelectBuilder, p *Pagination) (*Pagination, []T, error) {
	if p != nil {
		b = addPaginationQuery(b, p)
	}

	query, args, err := b.
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, nil, fmt.Errorf("build select with pagination: %w", err)
	}

	var res []T
	err = db.Get(&res, query, args...)
	if err != nil {
		return nil, nil, fmt.Errorf("select profile by subject id: %w", err)
	}

	if p == nil {
		return nil, res, nil
	}

	newP := *p
	if len(res) > p.Size {
		newP.Last.Key = res[p.Size-1]
	}

	return &newP, res[:p.Size], nil
}

func addPaginationQuery(b sq.SelectBuilder, p *Pagination) sq.SelectBuilder {
	if p.Sort != nil {
		order := AscSortLabel
		if !p.Sort.Asc {
			order = DescSortLabel
		}
		b = b.OrderBy(fmt.Sprintf("%s %s", p.Sort.Field, order))
	}

	if p.Sort == nil {
		p.Sort = &Sort{
			Field: "",
			Asc:   true,
		}
	}

	if p.Last != nil && p.Sort != nil {
		if !p.Sort.Asc {
			b = b.Where(sq.Lt{p.Last.Field: p.Last.Key})
		} else {
			b = b.Where(sq.Gt{p.Last.Field: p.Last.Key})
		}
	}

	b = b.Limit(uint64(p.Size) + 1)

	return b
}
