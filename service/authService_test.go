package service

import (
	"github.com/dgrijalva/jwt-go"
	realDomain "github.com/djedjethai/bankingAuth/domain"
	"github.com/djedjethai/bankingAuth/dto"
	"github.com/djedjethai/bankingAuth/errs"
	"github.com/djedjethai/bankingAuth/mocks/domain"
	"github.com/golang/mock/gomock"
	"testing"
	"time"
)

var mockRepo *domain.MockAuthRepository
var service *authService

// pb at testing type ???
func setup(t *testing.T) func() {

	ctrl := gomock.NewController(t)
	mockRepo = domain.NewMockAuthRepository(ctrl)

	// set the permissions, as they must be pass into the service
	permissions := realDomain.GetRolePermissions()

	service = NewService(mockRepo, permissions)

	return func() {
		service = nil
		defer ctrl.Finish()
	}
}

func Test_authService_login_return_an_err_if_domainReq_findById_return_err(t *testing.T) {
	tearDown := setup(t)
	defer tearDown()

	// Arrange
	lr := dto.LoginRequest{
		Username: "username",
		Password: "password",
	}

	mockRepo.EXPECT().FindBy("username", "password").Return(nil, errs.NewValidationError("invalid credential"))

	// Act
	_, err := service.Login(lr)

	// Assert
	if err == nil {
		t.Error("While testing authServ login should return an err if creadential does not exist in db")
	}
}

func Test_authService_login_return_an_err_if_domainReq_GenerateAndSaveToken_return_err(t *testing.T) {
	tearDown := setup(t)
	defer tearDown()

	// Arrange
	lr := dto.LoginRequest{
		Username: "2001",
		Password: "password",
	}

	// set login
	login := realDomain.Login{
		Username: "2001",
		Role:     "user",
	}

	mockRepo.EXPECT().FindBy("2001", "password").Return(&login, nil)

	// set token
	rtc := realDomain.AccessTokenClaims{
		CustomerId: "2001",
		Username:   "2001",
		Role:       "user",
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(realDomain.REFRESH_TOKEN_DURATION).Unix(),
		},
	}

	// PB HERERERERERE
	at := realDomain.NewAuthToken(rtc)
	tkn, _ := at.Token.SignedString([]byte(realDomain.HMAC_SAMPLE_SECRET))

	mockRepo.EXPECT().GenerateAndSaveRefreshTokenToStore(tkn).Return("", errs.NewInternalServerError("unexpected server error"))

	// Act
	_, err := service.Login(lr)

	// Assert
	if err == nil {
		t.Error("While testing authService Login should return an err if GenerateAndSaveToken return an err")
	}
}

// func Test_authService_login_return_dtoLoginResp(t *testing.T){}
