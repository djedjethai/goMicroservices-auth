package service

import (
	"github.com/djedjethai/bankingAuth/domain"
	"github.com/djedjethai/bankingAuth/dto"
	"github.com/djedjethai/bankingAuth/errs"
)

type AuthService interface {
	Login(dto.LoginRequest) (*dto.LoginResponse, *errs.AppError)
}

type Domain interface {
	FindBy(string, string) (*Loggin, *errs.AppError)
	GenerateAndSaveRefreshTokenToStore(AuthToken) (string, *errs.AppError)
}

type authService struct {
	repo Domain
}

func NewService(db Domain) *authService {
	return &authService{db}
}

func (s authService) Login(lr dto.LoginRequest) (*dto.LoginResponse, *errs.AppErr) {
	var login *domain.Login
	var appErr *errs.AppError

	login, appError = s.repo.FindBy(lr.Username, lr.Password)
	if appErr != nil {
		return nil, appErr
	}

	claims := login.ClaimsForAccessToken()
	authToken := domain.NewAuthToken(claims)

	var accessToken, refreshToken string
	if accessToken, appErr = authToken.NewAccessToken(); appErr != nil {
		return nil, appErr
	}

	if refreshToken, appErr = s.repo.GenerateAndSaveRefreshTokenToStore(authToken); err != nil {
		return nil, appErr
	}

	return &dto.LoginResponse{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}
