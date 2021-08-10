package service

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/djedjethai/bankingAuth/domain"
	"github.com/djedjethai/bankingAuth/dto"
	"github.com/djedjethai/bankingAuth/errs"
	"github.com/djedjethai/bankingAuth/logger"
)

type AuthService interface {
	Login(dto.LoginRequest) (*dto.LoginResponse, *errs.AppError)
	Verify(map[string]string) *errs.AppError
}

type authService struct {
	repo            domain.AuthRepository
	rolePermissions domain.RolePermissions
}

func NewService(db domain.AuthRepository, permissions domain.RolePermissions) *authService {
	return &authService{db, permissions}
}

func (s *authService) Verify(urlParams map[string]string) *errs.AppError {
	// convert the string token to JWT struct
	if jwtToken, err := jwtTokenFromString(urlParams["token"]); err != nil {
		return errs.NewAuthorizationError(err.Error())
	} else {

		// cerify the expire time and signature of the token
		if jwtToken.Valid {
			// type cast the token claims to jwt.MapClaims
			claims := jwtToken.Claims.(*domain.AccessTokenClaims)

			/* if role if user, then check if the account_id and customer_id
			   coming in the url belongs to the same token */
			if claims.IsUserRole() {
				if !claims.IsRequestVerifiedWithTokenClaims(urlParams) {
					return errs.NewAuthorizationError("request not verified with the token")
				}
			}
			// verify of the role is authorized to use the route
			isAuthorized := s.rolePermissions.IsAuthorizedFor(claims.Role, urlParams["routeName"])
			if !isAuthorized {
				return errs.NewAuthorizationError(fmt.Sprintf("%s role is not authorized", claims.Role))
			}

			return nil
		} else {
			return errs.NewAuthorizationError("Invalid token")
		}
	}

}

func (s *authService) Login(lr dto.LoginRequest) (*dto.LoginResponse, *errs.AppError) {
	var login *domain.Login
	var appErr *errs.AppError

	login, appErr = s.repo.FindBy(lr.Username, lr.Password)
	if appErr != nil {
		return nil, appErr
	}
	claims := login.ClaimsForAccessToken()
	authToken := domain.NewAuthToken(claims)

	var accessToken, refreshToken string
	if accessToken, appErr = authToken.NewAccessToken(); appErr != nil {
		return nil, appErr
	}

	if refreshToken, appErr = s.repo.GenerateAndSaveRefreshTokenToStore(authToken); appErr != nil {
		return nil, appErr
	}

	return &dto.LoginResponse{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}

func jwtTokenFromString(tokenString string) (*jwt.Token, error) {
	token, err := jwt.ParseWithClaims(tokenString, &domain.AccessTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(domain.HMAC_SAMPLE_SECRET), nil
	})
	if err != nil {
		logger.Error("Error while parsing token: " + err.Error())
		return nil, err
	}

	return token, nil
}
