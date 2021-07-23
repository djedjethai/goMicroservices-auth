package service

import (
	"github.com/djedjethai/bankingAuth/dto"
	"github.com/djedjethai/bankingAuth/errs"
)

type Service interface {
}

type Domain interface {
}

type service struct {
	db Domain
}

func NewService(db Domain) *service {
	return &service{db}
}

func (s Service) Login(dto.LoginRequest) (*dto.LoginResponse, *errs.AppErr) {

}
