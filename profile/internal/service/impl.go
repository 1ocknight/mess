package service

import "github.com/TATAROmangol/mess/profile/internal/storage"

type IMPL struct {
	s *storage.Service
}

func New(s *storage.Service) *IMPL {
	return &IMPL{
		s: s,
	}
}

