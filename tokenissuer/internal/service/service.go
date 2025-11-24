package service

import "tokenissuer/internal/adapter/identifier"

type Service interface {
	TokenService() Token
	VerifyService() Verify
}

type ServiceImpl struct {
	Token
	Verify
}

func NewServiceImpl(iden identifier.Service) *ServiceImpl {
	return &ServiceImpl{
		Token:  NewTokenImpl(iden),
		Verify: NewVerifyImpl(iden),
	}
}
