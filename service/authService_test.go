package service

import (
	// 	"github.com/dgrijalva/jwt-go"
	// "fmt"
	"github.com/dgrijalva/jwt-go"
	realDomain "github.com/djedjethai/bankingAuth/domain"
	"github.com/djedjethai/bankingAuth/dto"
	"github.com/djedjethai/bankingAuth/mocks/domain"
	"github.com/djedjethai/bankingLib/errs"
	"github.com/golang/mock/gomock"
	"net/http"
	"testing"
	"time"
)

var mockRepo *domain.MockAuthRepository
var service *authService

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

// ============= REFRESH ==============

func Test_authService_referesh_return_an_err_if_token_is_not_expired_yet(t *testing.T) {
	tearDown := setup(t)
	defer tearDown()

	// Arrange
	claim := realDomain.AccessTokenClaims{
		CustomerId: "2001",
		Username:   "2001",
		Role:       "user",
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(realDomain.ACCESS_TOKEN_DURATION).Unix(),
		},
	}
	authToken := realDomain.NewAuthToken(claim)
	token, _ := authToken.NewAccessToken()

	refreshedTokenReq := dto.RefreshTokenRequest{AccessToken: token}

	// Act
	_, err := service.Refresh(refreshedTokenReq)

	// Assert
	if err.Message != "cannot generate a new token before the current one expire" {
		t.Error("While testing authService Refresh should return an err if token is still not expire")
	}

}

func Test_authService_referesh_return_an_err_if_token_is_invalid(t *testing.T) {
	tearDown := setup(t)
	defer tearDown()

	// Arrange
	rt := dto.RefreshTokenRequest{AccessToken: "invalidToken"}

	// Act
	_, err := service.Refresh(rt)

	// Assert
	if err.Message != "invalide token" {
		t.Error("While testing authService refresh() should return an err if token is invalid")
	}
}

func Test_authService_referesh_return_an_err_if_domain_newAccessTokenFromRefreshToken_return_an_err(t *testing.T) {
	tearDown := setup(t)
	defer tearDown()

	// Arrange
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0b2tlbl90eXBlIjoicmVmcmVzaF90b2tlbiIsImN1c3RvbWVyX2lkIjoiIiwiYWNjb3VudHMiOm51bGwsInVuIjoiYWRtaW4iLCJyb2xlIjoiYWRtaW4iLCJleHAiOjE2MjkzODE4MjB9.dtsik_uSKfoduArFg0ZuneApz9IfNN0rOL1rS-ByuM8"

	rt := dto.RefreshTokenRequest{AccessToken: token, RefreshToken: token}

	mockRepo.EXPECT().RefreshTokenExists(rt.RefreshToken).Return(errs.NewNotFoundError("not found"))

	// Act
	_, err := service.Refresh(rt)

	// Assert
	if err.Code != http.StatusNotFound {
		t.Error("While testing authService refresh() should return an err if domain.newAccessTokenFromRefreshToken() return one")
	}
}

// A FINIR(return): ahhhh: &{401 invalid or expired refresh token}
// If i generate a new one it's still valid and return an err
// If expire return an err, so how should be the token to be renew ???
// func Test_authService_referesh_should_not_return_an_err(t *testing.T) {
// 	tearDown := setup(t)
// 	defer tearDown()
//
// 	// Arrange
// 	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0b2tlbl90eXBlIjoicmVmcmVzaF90b2tlbiIsImN1c3RvbWVyX2lkIjoiIiwiYWNjb3VudHMiOm51bGwsInVuIjoiYWRtaW4iLCJyb2xlIjoiYWRtaW4iLCJleHAiOjE2MjkzODE4MjB9.dtsik_uSKfoduArFg0ZuneApz9IfNN0rOL1rS-ByuM8"
//
// 	rt := dto.RefreshTokenRequest{AccessToken: token, RefreshToken: token}
//
// 	mockRepo.EXPECT().RefreshTokenExists(rt.RefreshToken).Return(nil)
//
// 	// Act
// 	ret, err := service.Refresh(rt)
//
// 	fmt.Printf("ahhhh: %v\n", err)
// 	fmt.Printf("ahhhh222: %v\n", ret)
//
// 	// Assert
// 	if err != nil {
// 		t.Error("While testing authService refresh() sould not return any err")
// 	}
// }

// ============= VERIFY ================
func Test_authService_verify_return_an_err_if_the_map_token_key_is_incorrect(t *testing.T) {
	// Arrange
	tearDown := setup(t)
	defer tearDown()

	urlParams := make(map[string]string)
	urlParams["token"] = "wrongToken"

	// Act
	err := service.Verify(urlParams)

	// Assert
	if err.Message != "token contains an invalid number of segments" {
		t.Error("While testing authService, Verify() should return an err if token is invalid ")
	}
}

func Test_authService_verify_return_an_err_if_urlParamsToken_is_expired(t *testing.T) {
	// Arrange
	tearDown := setup(t)
	defer tearDown()

	urlParams := make(map[string]string)

	// expired token
	urlParams["token"] = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0b2tlbl90eXBlIjoicmVmcmVzaF90b2tlbiIsImN1c3RvbWVyX2lkIjoiIiwiYWNjb3VudHMiOm51bGwsInVuIjoiYWRtaW4iLCJyb2xlIjoiYWRtaW4iLCJleHAiOjE2MjkzODE4MjB9.dtsik_uSKfoduArFg0ZuneApz9IfNN0rOL1rS-ByuM8"

	// Act
	err := service.Verify(urlParams)

	// Assert
	if err.Code != http.StatusForbidden {
		t.Error("While testing authService, Verify() should return an err if token is invalid ")
	}
}

func Test_authService_verify_return_an_err_if_userRole_and_IsRequestVerifiedWithTokenClaims_is_false(t *testing.T) {
	// Arrange
	tearDown := setup(t)
	defer tearDown()

	urlParams := make(map[string]string)

	// create a valid token
	claim := realDomain.AccessTokenClaims{
		CustomerId: "2001",
		Username:   "2001",
		Role:       "wronguser",
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(realDomain.ACCESS_TOKEN_DURATION).Unix(),
		},
	}
	authToken := realDomain.NewAuthToken(claim)
	token, _ := authToken.NewAccessToken()

	urlParams["token"] = token

	// Act
	err := service.Verify(urlParams)

	// Assert
	if err.Message != "wronguser role is not authorized" {
		t.Error("While testing authService, if wrong userRole, Verify() should return an err")
	}

}

func Test_authService_verify_return_an_err_if_token_is_valid_but_is_not_authorized(t *testing.T) {
	// Arrange
	tearDown := setup(t)
	defer tearDown()

	urlParams := make(map[string]string)

	// create a valid token
	claim := realDomain.AccessTokenClaims{
		CustomerId: "2001",
		Username:   "2001",
		Role:       "user",
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(realDomain.ACCESS_TOKEN_DURATION).Unix(),
		},
	}
	authToken := realDomain.NewAuthToken(claim)
	token, _ := authToken.NewAccessToken()

	urlParams["token"] = token
	urlParams["customer_id"] = "wrong"

	// Act
	err := service.Verify(urlParams)

	// Assert
	if err.Message != "request not verified with the token" {
		t.Error("While testing authService, if wrong userRole, Verify() should return an err")
	}
}

func Test_authService_verify_does_not_return_any_err(t *testing.T) {
	// Arrange
	tearDown := setup(t)
	defer tearDown()

	urlParams := make(map[string]string)

	// create a valid token
	claim := realDomain.AccessTokenClaims{
		CustomerId: "2001",
		Username:   "2001",
		Role:       "admin",
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(realDomain.ACCESS_TOKEN_DURATION).Unix(),
		},
	}
	authToken := realDomain.NewAuthToken(claim)
	token, _ := authToken.NewAccessToken()

	urlParams["token"] = token
	urlParams["customer_id"] = "2001"
	urlParams["routeName"] = "GetCustomer"

	// Act
	err := service.Verify(urlParams)

	// Assert
	if err != nil {
		t.Error("While testing authService, Verify() should not return any err")
	}

}

// ============= SIGNUP ================
// test 2 methods ValidUsername and ValidNameDobCityZip
func Test_authService_signup_should_return_an_err_if_checkUserInput_is_incorrect(t *testing.T) {
	tearDown := setup(t)
	defer tearDown()

	// Arrange
	sr1 := dto.SignupRequest{
		Name:         "myname",
		DateOfBirth:  "1972-05-03",
		City:         "mycity",
		ZipCode:      "10000",
		Username:     "my",
		Password:     "password",
		PasswordConf: "password",
	}

	sr2 := dto.SignupRequest{
		Name:         "myname",
		DateOfBirth:  "1972-05-03",
		City:         "m",
		ZipCode:      "10000",
		Username:     "mypassword",
		Password:     "password",
		PasswordConf: "password",
	}

	// Act
	_, err1 := service.Signup(sr1)
	_, err2 := service.Signup(sr2)

	// Assert
	if err1 == nil {
		t.Error("While testing authService signup should return an err if ValidNewUser return an err")
	}

	if err2 == nil {
		t.Error("While testing authService signup should return an err if ValidNameDobCityZip return an err")
	}

}

func Test_authService_signup_return_an_err_if_username_exist(t *testing.T) {
	tearDown := setup(t)
	defer tearDown()

	// Arrange
	sr := dto.SignupRequest{
		Name:         "jerome",
		DateOfBirth:  "1972-05-03",
		City:         "madrid",
		ZipCode:      "10000",
		Username:     "myusername",
		Password:     "password",
		PasswordConf: "password",
	}

	mockRepo.EXPECT().IsUsernameExist("myusername").Return(true, nil)

	// Act
	_, err := service.Signup(sr)

	// Assert
	if err.Code != http.StatusUnprocessableEntity {
		t.Error("While testing authService signup should return an err if username already exist")
	}
}

func Test_authService_signup_return_an_err_if_IfUsernameExist_query_return_an_err(t *testing.T) {
	tearDown := setup(t)
	defer tearDown()

	// Arrange
	sr := dto.SignupRequest{
		Name:         "jerome",
		DateOfBirth:  "1972-05-03",
		City:         "madrid",
		ZipCode:      "10000",
		Username:     "myusername",
		Password:     "password",
		PasswordConf: "password",
	}

	mockRepo.EXPECT().IsUsernameExist("myusername").Return(false, errs.NewInternalServerError("Unexpected database err"))

	// Act
	_, err := service.Signup(sr)

	// Assert
	if err.Code != http.StatusInternalServerError {
		t.Error("While testing authService signup should return an err if IsUserExist return an InternalServerError")
	}
}

func Test_authService_signup_return_an_err_if_query_CreateCustAndUser_return_an_err(t *testing.T) {
	tearDown := setup(t)
	defer tearDown()

	// Arrange
	sr := dto.SignupRequest{
		Name:         "jerome",
		DateOfBirth:  "1972-05-03",
		City:         "madrid",
		ZipCode:      "10000",
		Username:     "myusername",
		Password:     "password",
		PasswordConf: "password",
	}

	cd := realDomain.CustomerDomain{
		Name:        sr.Name,
		DateOfBirth: sr.DateOfBirth,
		City:        sr.City,
		ZipCode:     sr.ZipCode,
		Username:    sr.Username,
		Password:    sr.Password,
	}

	mockRepo.EXPECT().IsUsernameExist("myusername").Return(false, nil)

	mockRepo.EXPECT().CreateCustAndUser(cd).Return(nil, errs.NewInternalServerError("Unexpected db err"))

	// Act
	_, err := service.Signup(sr)

	// Assert
	if err.Code != http.StatusInternalServerError {
		t.Error("While testing authService signup should return an err if CreateCustAndUser query return an err")
	}
}

func Test_authService_signup_return_an_err_if_GenerateAndSaveRfereshTokenToStore_return_an_err(t *testing.T) {
	tearDown := setup(t)
	defer tearDown()

	// Arrange
	sr := dto.SignupRequest{
		Name:         "jerome",
		DateOfBirth:  "1972-05-03",
		City:         "madrid",
		ZipCode:      "10000",
		Username:     "myusername",
		Password:     "password",
		PasswordConf: "password",
	}

	cd := realDomain.CustomerDomain{
		Name:        sr.Name,
		DateOfBirth: sr.DateOfBirth,
		City:        sr.City,
		ZipCode:     sr.ZipCode,
		Username:    sr.Username,
		Password:    sr.Password,
	}

	// should provoc an err
	lgi := realDomain.Login{
		Username: "2001",
		Role:     "user",
	}

	mockRepo.EXPECT().IsUsernameExist("myusername").Return(false, nil)

	mockRepo.EXPECT().CreateCustAndUser(cd).Return(&lgi, nil)

	claims := lgi.ClaimsForAccessToken()
	authToken := realDomain.NewAuthToken(claims)

	mockRepo.EXPECT().GenerateAndSaveRefreshTokenToStore(authToken).Return("", errs.NewInternalServerError("unexpected db error"))

	// Act
	_, err := service.Signup(sr)

	// Assert
	if err.Code != http.StatusInternalServerError {
		t.Error("While testing authService signup should return an err if CreateCustAndUser query return an err")
	}
}

func Test_authService_signup_do_not_return_an_err_(t *testing.T) {
	tearDown := setup(t)
	defer tearDown()

	// Arrange
	sr := dto.SignupRequest{
		Name:         "jerome",
		DateOfBirth:  "1972-05-03",
		City:         "madrid",
		ZipCode:      "10000",
		Username:     "myusername",
		Password:     "password",
		PasswordConf: "password",
	}

	cd := realDomain.CustomerDomain{
		Name:        sr.Name,
		DateOfBirth: sr.DateOfBirth,
		City:        sr.City,
		ZipCode:     sr.ZipCode,
		Username:    sr.Username,
		Password:    sr.Password,
	}

	// should provoc an err
	lgi := realDomain.Login{
		Username: "2001",
		Role:     "user",
	}

	mockRepo.EXPECT().IsUsernameExist("myusername").Return(false, nil)

	mockRepo.EXPECT().CreateCustAndUser(cd).Return(&lgi, nil)

	claims := lgi.ClaimsForAccessToken()
	authToken := realDomain.NewAuthToken(claims)

	// generate a refreshed token to be return from GenerateAndSave...
	rft, _ := authToken.NewRefreshToken()

	mockRepo.EXPECT().GenerateAndSaveRefreshTokenToStore(authToken).Return(rft, nil)

	// Act
	_, err := service.Signup(sr)

	// Assert
	if err != nil {
		t.Error("While testing authService signup should not return an err")
	}
}

// ============= LOGIN ================
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

	claims := login.ClaimsForAccessToken()
	at := realDomain.NewAuthToken(claims)

	mockRepo.EXPECT().GenerateAndSaveRefreshTokenToStore(at).Return("", errs.NewInternalServerError("unexpected server error"))

	// Act
	_, err := service.Login(lr)

	// Assert
	if err == nil {
		t.Error("While testing authService Login should return an err if GenerateAndSaveToken return an err")
	}
}

func Test_authService_login_return_dtoLoginResp(t *testing.T) {
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

	claims := login.ClaimsForAccessToken()
	authToken := realDomain.NewAuthToken(claims)

	// set a refresh token
	refreshedToken, _ := authToken.NewRefreshToken()

	mockRepo.EXPECT().GenerateAndSaveRefreshTokenToStore(authToken).Return(refreshedToken, nil)

	// Act
	_, err := service.Login(lr)

	// Assert
	if err != nil {
		t.Error("While testing authService Login should not return an err")
	}
}
