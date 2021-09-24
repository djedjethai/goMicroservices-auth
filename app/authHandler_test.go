package app

import (
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"net/http"
	"net/http/httptest"
	"testing"
)

// NEED TO GET THE PACKAGE GOMOCK (for the mock to be use)
// go install github.com/golang/mock/mockgen@v1.6.0 (at this day)

var router *mux.Router
var mockService *serviceAuth
var ah authHandler

func setup(t *testing.T) func(){
	ctrl := gomock.NewController()

	mockService := service.NewMockAuthService(ctrl)

	ah = authHandler{mockService}

	return func() {
		router nil
		defer ctrl.Finish()
	}
}

// test login
// func Test_authHandler_login_method_should_return_httpErr_if_body_is_not_json(t *Testing.T){}
// func Test_authHandler_login_method_should_return_statusCodeErr_if_wrong_credential(t *Testing.T){}
// func Test_authHandler_login_method_should_return_statusCodeOK_if_wrong_credential(t *Testing.T){}
