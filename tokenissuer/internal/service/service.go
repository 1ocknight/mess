package service

import (
	"time"
	"tokenissuer/internal/adapter/identifier"
)

type Service interface {
	TokenService() Token
	VerifyService() Verify
}

type ServiceImpl struct {
	Token
	Verify
}

func NewServiceImpl(iden identifier.Service, jwksTTL time.Duration) *ServiceImpl {
	return &ServiceImpl{
		Token:  NewTokenImpl(iden),
		Verify: NewVerifyImpl(iden, jwksTTL),
	}
}
