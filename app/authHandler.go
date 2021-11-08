package app

import (
	"encoding/json"
	"fmt"
	"github.com/djedjethai/bankingAuth/dto"
	// "github.com/djedjethai/bankingAuth/errs"
	"github.com/djedjethai/bankingAuth/service"
	"github.com/djedjethai/bankingLib/logger"
	"net/http"
)

type authHandler struct {
	service service.AuthService
}

func (h authHandler) addCustomer(w http.ResponseWriter, r *http.Request) {
	var signupRequest dto.SignupRequest
	if err := json.NewDecoder(r.Body).Decode(&signupRequest); err != nil {
		logger.Error("error when new customer signup")
		w.WriteHeader(http.StatusBadRequest)
	} else {
		token, appErr := h.service.Signup(signupRequest)
		if appErr != nil {
			fmt.Println("grrree 2: ", appErr)
			writeResponse(w, appErr.Code, appErr.AsMessage())
		} else {
			writeResponse(w, http.StatusOK, token)
		}
	}
}

func (h authHandler) refresh(w http.ResponseWriter, r *http.Request) {
	var refreshRequest dto.RefreshTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&refreshRequest); err != nil {
		logger.Error("Error while decoding refresh token request" + err.Error())
		w.WriteHeader(http.StatusBadRequest)
	} else {
		token, appErr := h.service.Refresh(refreshRequest)
		if appErr != nil {
			writeResponse(w, appErr.Code, appErr.AsMessage())
		} else {
			writeResponse(w, http.StatusOK, *token)
		}
	}
}

func (h authHandler) verify(w http.ResponseWriter, r *http.Request) {
	urlParams := make(map[string]string)

	// converting from query to map type
	for k := range r.URL.Query() {
		urlParams[k] = r.URL.Query().Get(k)
	}

	if urlParams["token"] != "" {
		appErr := h.service.Verify(urlParams)
		if appErr != nil {
			writeResponse(w, appErr.Code, notAuthorizedResponse(appErr.Message))
		} else {
			writeResponse(w, http.StatusOK, authorizedResponse())
		}
	} else {
		writeResponse(w, http.StatusForbidden, notAuthorizedResponse("token missing"))
	}
}

func (h authHandler) login(w http.ResponseWriter, r *http.Request) {
	var loginReq dto.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&loginReq); err != nil {
		logger.Error("Error when decode login request" + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	token, appErr := h.service.Login(loginReq)
	if appErr != nil {
		logger.Error("In authHandler service.Login")
		writeResponse(w, appErr.Code, appErr.AsMessage())
		return
	}

	writeResponse(w, http.StatusOK, *token)
}

func notAuthorizedResponse(msg string) map[string]interface{} {
	return map[string]interface{}{
		"isAuthorized": false,
		"message":      msg,
	}
}

func authorizedResponse() map[string]bool {
	return map[string]bool{"IsAuthorized": true}
}

func writeResponse(w http.ResponseWriter, code int, data interface{}) {
	w.Header().Add("Content-type", "application/json")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		logger.Error("Error when writter encode datas")
		panic(err)
	}
}
