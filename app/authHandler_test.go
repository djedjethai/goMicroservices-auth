package app

import (
	"github.com/djedjethai/bankingAuth/dto"
	"github.com/djedjethai/bankingAuth/errs"
	"github.com/djedjethai/bankingAuth/mocks/service"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// NEED TO GET THE PACKAGE GOMOCK (for the mock to be use)
// go install github.com/golang/mock/mockgen@v1.6.0 (at this day)
// go get: added github.com/golang/mock v1.6.0

var router *mux.Router
var mockService *service.MockAuthService
var ah authHandler

func setup(t *testing.T) func() {
	ctrl := gomock.NewController(t)

	mockService = service.NewMockAuthService(ctrl)

	ah = authHandler{mockService}

	router = mux.NewRouter()

	return func() {
		router = nil
		defer ctrl.Finish()
	}
}

func Test_authHandler_verify_should_return_err_if_url_does_not_hold_any_params(t *testing.T) {
	tearDown := setup(t)
	defer tearDown()

	// Arrange
	router.HandleFunc("/auth/verify", ah.verify)
	request, _ := http.NewRequest(http.MethodGet, "/auth/verify", nil)

	// Act
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	// Assert
	if recorder.Code != http.StatusForbidden {
		t.Error("While testing authHandler verify should return an err if no params in url")
	}
}

func Test_authHandler_verify_should_return_err_if_service_return_err(t *testing.T) {
	tearDown := setup(t)
	defer tearDown()

	// Arrange
	// set serviceInput
	serviceInput := make(map[string]string)
	serviceInput["token"] = "wrongToken"

	mockService.EXPECT().Verify(serviceInput).Return(errs.NewAuthorizationError("request not verified with the token"))
	router.HandleFunc("/auth/verify", ah.verify)
	request, _ := http.NewRequest(http.MethodGet, "/auth/verify?token=wrongToken", nil)

	// Act
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	// Assert
	if recorder.Code != http.StatusForbidden {
		t.Error("While testing authService verify should return an err if verifyService return an err")
	}
}

func Test_authHandler_verify_should_return_statusCodeOK_if_service_return_nil(t *testing.T) {
	tearDown := setup(t)
	defer tearDown()

	// Arrange
	serviceInput := make(map[string]string)
	serviceInput["token"] = "theToken"

	router.HandleFunc("/auth/verify", ah.verify)
	request, _ := http.NewRequest(http.MethodGet, "/auth/verify?token=theToken", nil)

	mockService.EXPECT().Verify(serviceInput).Return(nil)

	// Act
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	// Assert
	if recorder.Code != http.StatusOK {
		t.Error("While testing authHandler verify token should return statusCode ok")
	}

}

func Test_authHandler_refreshToken_should_return_err_if_body_is_not_json(t *testing.T) {
	tearDown := setup(t)
	defer tearDown()

	// Arrange
	router.HandleFunc("/auth/refresh", ah.refresh)
	request, _ := http.NewRequest(http.MethodPost, "/auth/refresh", strings.NewReader(""))

	// Act
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	// Assert
	if recorder.Code != http.StatusBadRequest {
		t.Error("While testing authHandler refreshToken should return an err when body is not json")
	}
}

func Test_authHandler_refreshToken_should_return_err_if_service_return_an_err(t *testing.T) {
	tearDown := setup(t)
	defer tearDown()

	// Arrange
	// set jsonInput
	jsonInput := `{
		"access_token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9",
		"refresh_token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9"
	}`
	// set serviceInput
	serviceInput := dto.RefreshTokenRequest{
		AccessToken:  "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9",
		RefreshToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9",
	}

	router.HandleFunc("/auth/refresh", ah.refresh)
	request, _ := http.NewRequest(http.MethodPost, "/auth/refresh", strings.NewReader(jsonInput))
	mockService.EXPECT().Refresh(serviceInput).Return(nil, errs.NewAuthenticationError("invalid token"))

	// Act
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	// Assert
	if recorder.Code != http.StatusUnauthorized {
		t.Error("While testing authHandler RefreshToken should return an err if token is invalid")
	}
}

func Test_authHandler_refreshToken_should_return_StatusOK_if_service_refreshed_token(t *testing.T) {
	tearDown := setup(t)
	defer tearDown()

	// Arrange
	jsonInput := `{
		"access_token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9",
		"refresh_token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9"
	}`
	// set serviceInput
	serviceInput := dto.RefreshTokenRequest{
		AccessToken:  "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9",
		RefreshToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9",
	}

	// set serviceOutput
	serviceOutput := dto.LoginResponse{
		AccessToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0b2tlbl90eXBlIjoicmVmcmVzaF90b2tlbiIsImN1c3RvbWVyX2lkIjoiIiwiYWNjb3VudHMiOm51bGwsInVuIjoiYWRtaW4iLCJyb2xlIjoiYWRtaW4iLCJleHAiOjE2MjkzODE4MjB9.dtsik_uSKfoduArFg0ZuneApz9IfNN0rOL1rS-ByuM8",
	}

	mockService.EXPECT().Refresh(serviceInput).Return(&serviceOutput, nil)
	router.HandleFunc("/auth/refresh", ah.refresh)
	request, _ := http.NewRequest(http.MethodPost, "/auth/refresh", strings.NewReader(jsonInput))
	// Act
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	// Assert
	if recorder.Code != http.StatusOK {
		t.Error("While testing authHandler refreshToken should return statusCode ok")
	}
}

func Test_authHandler_addCustomer_should_return_err_if_input_is_not_json(t *testing.T) {
	tearDown := setup(t)
	defer tearDown()

	// Arrange
	router.HandleFunc("/auth/register", ah.addCustomer)
	request, _ := http.NewRequest(http.MethodPost, "/auth/register", strings.NewReader(""))

	// Act
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	// Assert
	if recorder.Code != http.StatusBadRequest {
		t.Error("While testing authHandler addCustomer body request must be json")
	}
}

func Test_authHandler_addCustomer_should_return_err_if_service_return_an_err(t *testing.T) {
	tearDown := setup(t)
	defer tearDown()

	// Arrange
	jsonInput := `{
    		"name":"jerome",
    		"date_of_birth":"1975-05-03",
    		"city":"bangkok",
    		"zip_code":"10000",
    		"username":"jerome",
    		"password":"password",
    		"password_conf":"password"
	}`

	serviceInput := dto.SignupRequest{
		Name:         "jerome",
		DateOfBirth:  "1975-05-03",
		City:         "bangkok",
		ZipCode:      "10000",
		Username:     "jerome",
		Password:     "password",
		PasswordConf: "password",
	}

	serviceOutput := dto.LoginResponse{AccessToken: "qwerty", RefreshToken: "qwerty"}

	mockService.EXPECT().Signup(serviceInput).Return(&serviceOutput, nil)

	router.HandleFunc("/auth/register", ah.addCustomer)
	request, _ := http.NewRequest(http.MethodPost, "/auth/register", strings.NewReader(jsonInput))
	// Act
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	// Assert
	if recorder.Code != http.StatusOK {
		t.Error("While testing authHandler addCustomer with correct credential should return statusOK")
	}

}

func Test_authHandler_addCustomer_should_return_err_if_service_return_err(t *testing.T) {
	tearDown := setup(t)
	defer tearDown()

	// Arrange
	jsonInput := `{
    		"name":"jerome",
    		"date_of_birth":"1975-05-03",
    		"city":"bangkok",
    		"zip_code":"10000",
    		"username":"jerome",
    		"password":"password",
    		"password_conf":"passw"
	}`

	serviceInput := dto.SignupRequest{
		Name:         "jerome",
		DateOfBirth:  "1975-05-03",
		City:         "bangkok",
		ZipCode:      "10000",
		Username:     "jerome",
		Password:     "password",
		PasswordConf: "passw",
	}

	mockService.EXPECT().Signup(serviceInput).Return(nil, errs.NewValidationError("wrong input"))

	router.HandleFunc("/auth/register", ah.addCustomer)
	request, _ := http.NewRequest(http.MethodPost, "/auth/register", strings.NewReader(jsonInput))
	// Act
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	// Assert
	if recorder.Code != http.StatusUnprocessableEntity {
		t.Error("While testing authService, addCustomer should return an err if service return an err")
	}
}

// test Login()
func Test_authHandler_login_method_should_return_statusCodeOK_if_correct_credential(t *testing.T) {
	tearDown := setup(t)
	defer tearDown()

	// Arrange
	jsonInput := `{"username":"admin", "password":"admin123"}`

	serviceInput := dto.LoginRequest{Username: "admin", Password: "admin123"}

	serviceOutput := dto.LoginResponse{AccessToken: "qwerty", RefreshToken: "qwerty"}

	mockService.EXPECT().Login(serviceInput).Return(&serviceOutput, nil)

	router.HandleFunc("/auth/login", ah.login)
	request, _ := http.NewRequest(http.MethodPost, "/auth/login", strings.NewReader(jsonInput))

	// Act
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	// Assert
	if recorder.Code != http.StatusOK {
		t.Error("While testing authHandler login with corrwct credential should return statusOK")
	}
}

func Test_authHandler_login_method_should_return_statusCodeErr_if_wrong_credential(t *testing.T) {
	tearDown := setup(t)
	defer tearDown()

	// Arrange
	jsonInput := `{"username":"admin", "password":"wrongpassword"}`

	loginServiceInput := dto.LoginRequest{
		Username: "admin",
		Password: "wrongpassword",
	}

	mockService.EXPECT().Login(loginServiceInput).Return(nil, errs.NewValidationError("invalid credential"))

	router.HandleFunc("/auth/login", ah.login)
	request, _ := http.NewRequest(http.MethodPost, "/auth/login", strings.NewReader(jsonInput))

	// Act
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	// Assert
	if recorder.Code != http.StatusUnprocessableEntity {
		t.Error("While testing authHandler login with wrong creadential should return an err")
	}
}

// test login
func Test_authHandler_login_method_should_return_httpErr_if_body_is_not_json(t *testing.T) {
	tearDown := setup(t)
	defer tearDown()

	// request without json input
	// Arrange
	router.HandleFunc("/auth/login", ah.login)
	request, _ := http.NewRequest(http.MethodPost, "/auth/login", strings.NewReader(""))

	// Act
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	// Assert
	if recorder.Code != http.StatusBadRequest {
		t.Error("While testing authService login handler should return an err if input is not json format")
	}

}
