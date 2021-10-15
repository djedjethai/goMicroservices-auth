package service

import (
	// 	"github.com/dgrijalva/jwt-go"
	realDomain "github.com/djedjethai/bankingAuth/domain"
	"github.com/djedjethai/bankingAuth/dto"
	"github.com/djedjethai/bankingAuth/errs"
	"github.com/djedjethai/bankingAuth/mocks/domain"
	"github.com/golang/mock/gomock"
	"net/http"
	"testing"
	// 	"time"
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

// func Test_authService_signup_return_an_err_if_query_IsUserName_return_an_err(t *testing.T){
// 	tearDown := setup(t)
// 	defer tearDown()
//
// 	// Arrange
// }

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

// A FINIR
// func Test_authService_signup_return_an_err_if_authTokenNewAccessToken_return_an_err(t *testing.T) {
// 	tearDown := setup(t)
// 	defer tearDown()
//
// 	// Arrange
// 	sr := dto.SignupRequest{
// 		Name:         "jerome",
// 		DateOfBirth:  "1972-05-03",
// 		City:         "madrid",
// 		ZipCode:      "10000",
// 		Username:     "myusername",
// 		Password:     "password",
// 		PasswordConf: "password",
// 	}
//
// 	cd := realDomain.CustomerDomain{
// 		Name:        sr.Name,
// 		DateOfBirth: sr.DateOfBirth,
// 		City:        sr.City,
// 		ZipCode:     sr.ZipCode,
// 		Username:    sr.Username,
// 		Password:    sr.Password,
// 	}
//
// 	// should provoc an err
// 	lgi := realDomain.Login{
// 		Username: "2001",
// 		Role: "user",
// 	}
//
// 	mockRepo.EXPECT().IsUsernameExist("myusername").Return(false, nil)
//
// 	mockRepo.EXPECT().CreateCustAndUser(cd).Return(&lgi, nil)
//
// 	// claims := lgi.ClaimsForAccessToken()
// 	// authToken := realDomain.NewAuthToken(claims)
//
// 	// NEED TO BE MOCK TO TEST IF IN CASE OF ERR THE FUNC RETURN THE ERR
// 	// _, err1 := authToken.NewAccessToken()
//
// 	// Act
// 	_, err := service.Signup(sr)
//
// 	// Assert
// 	if err.Code != http.StatusInternalServerError {
// 		t.Error("While testing authService signup should return an err if CreateCustAndUser query return an err")
// 	}
//
// }

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

// TO FINISHHH
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

	// TO FINISH HERE RETURN NO ERRRRRRRR
	mockRepo.EXPECT().GenerateAndSaveRefreshTokenToStore(authToken).Return()

	// Act
	_, err := service.Signup(sr)

	// Assert
	if err.Code != http.StatusInternalServerError {
		t.Error("While testing authService signup should return an err if CreateCustAndUser query return an err")
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
