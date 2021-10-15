package service

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/djedjethai/bankingAuth/domain"
	"github.com/djedjethai/bankingAuth/dto"
	"github.com/djedjethai/bankingAuth/errs"
	"github.com/djedjethai/bankingAuth/logger"
	// "time"
)

//go:generate mockgen -destination=../mocks/service/mockAuthService.go -package=service github.com/djedjethai/bankingAuth/service AuthService
type AuthService interface {
	Login(dto.LoginRequest) (*dto.LoginResponse, *errs.AppError)
	Verify(map[string]string) *errs.AppError
	Refresh(dto.RefreshTokenRequest) (*dto.LoginResponse, *errs.AppError)
	Signup(dto.SignupRequest) (*dto.LoginResponse, *errs.AppError)
}

type authService struct {
	repo            domain.AuthRepository
	rolePermissions domain.RolePermissions
}

func NewService(db domain.AuthRepository, permissions domain.RolePermissions) *authService {
	return &authService{db, permissions}
}

func (s *authService) Refresh(request dto.RefreshTokenRequest) (*dto.LoginResponse, *errs.AppError) {
	if vErr := request.IsAccessTokenValid(); vErr != nil {
		if vErr.Errors == jwt.ValidationErrorExpired {
			// continue with the refresh token functionality
			var appErr *errs.AppError

			if appErr = s.repo.RefreshTokenExists(request.RefreshToken); appErr != nil {
				return nil, appErr
			}

			// generate a access token from refresh token
			var accessToken string
			if accessToken, appErr = domain.NewAccessTokenFromRefreshToken(request.RefreshToken); appErr != nil {
				return nil, appErr
			}
			return &dto.LoginResponse{AccessToken: accessToken}, nil
		}
		return nil, errs.NewAuthenticationError("invalide token")
	}
	return nil, errs.NewAuthenticationError("cannot generate a new token before the current one expire")
}

func (s *authService) Verify(urlParams map[string]string) *errs.AppError {

	logger.Info("get token for identification in auth svc" + urlParams["token"])

	// convert the string token to JWT struct
	if jwtToken, err := jwtTokenFromString(urlParams["token"]); err != nil {
		return errs.NewAuthorizationError(err.Error())
	} else {

		// cerify the expire time and signature of the token
		if jwtToken.Valid {
			// type cast the token claims to jwt.MapClaims
			claims := jwtToken.Claims.(*domain.AccessTokenClaims)

			logger.Info("Auth Verify after jwtToken.Valid: " + fmt.Sprintf("%v", claims))
			/* if role if user, then check if the account_id and customer_id
			   coming in the url belongs to the same token */
			if claims.IsUserRole() {
				logger.Info("In claims.IsUserRole")
				if !claims.IsRequestVerifiedWithTokenClaims(urlParams) {
					return errs.NewAuthorizationError("request not verified with the token")
				}
			}
			// verify of the role is authorized to use the route
			isAuthorized := s.rolePermissions.IsAuthorizedFor(claims.Role, urlParams["routeName"])
			logger.Info("check if authorized after verif rolePermission: " + fmt.Sprintf("%v", isAuthorized))
			if !isAuthorized {
				return errs.NewAuthorizationError(fmt.Sprintf("%s role is not authorized", claims.Role))
			}

			return nil
		} else {
			return errs.NewAuthorizationError("Invalid token")
		}
	}

}

// func (s *authService) signup(sr dto.SignupRequest) (*dto.LoginResponse, *errs.AppError) {
func (s *authService) Signup(sr dto.SignupRequest) (*dto.LoginResponse, *errs.AppError) {

	// check user input
	if err := sr.ValidNewUser(); err != nil {
		return nil, err
	}
	if err := sr.ValidNameDobCityZip(); err != nil {
		return nil, err
	}

	// check if username is avaible, I SHOULD CHECK OTHER CREDENTIAL, but i don't
	usernameExist, err := s.repo.IsUsernameExist(sr.Username)
	if err != nil {
		// return nil, err
		return nil, err
	}
	if usernameExist {
		// return nil, errs.NewValidationError("Username is not available")
		return nil, errs.NewValidationError("Invalid username")
	}

	// create customer
	custDom := domain.CustomerDomain{
		Name:        sr.Name,
		DateOfBirth: sr.DateOfBirth,
		City:        sr.City,
		ZipCode:     sr.ZipCode,
		Username:    sr.Username,
		Password:    sr.Password,
	}

	// create User(username, need the id from previous req)
	// login, err := s.repo.CreateCustAndUser(custDom)
	var login *domain.Login
	login, err = s.repo.CreateCustAndUser(custDom)
	if err != nil {
		return nil, err
	}

	claims := login.ClaimsForAccessToken()
	authToken := domain.NewAuthToken(claims)

	var accessToken, refreshToken string
	if accessToken, err = authToken.NewAccessToken(); err != nil {
		return nil, err
	}

	if refreshToken, err = s.repo.GenerateAndSaveRefreshTokenToStore(authToken); err != nil {
		return nil, err
	}

	return &dto.LoginResponse{AccessToken: accessToken, RefreshToken: refreshToken}, nil
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
		logger.Info("parsing the token in auth verify")
		return []byte(domain.HMAC_SAMPLE_SECRET), nil
	})
	if err != nil {
		logger.Error("Error while parsing token: " + err.Error())
		return nil, err
	}

	logger.Info("auth verify token done: " + fmt.Sprintf("%v", token))
	return token, nil
}
