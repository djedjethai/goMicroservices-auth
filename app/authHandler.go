package app

import (
	"encoding/json"
	"github.com/djedjethai/bankingAuth/dto"
	"github.com/djedjethai/bankingAuth/errs"
	"github.com/djedjethai/bankingAuth/logger"
	"github.com/djedjethai/bankingAuth/service"
	"net/http"
)

type authHandler struct {
	service service.AuthService
}

func writeResponse(w http.ResponseWriter, code int, data interface{}) {
	w.Header().Add("Content-type", "application/json")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		logger.Error("Error when writter encode datas")
		panic(err)
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
		writeResponse(w, appErr.Code, appErr.AsMessage())
	}

	writeResponse(w, http.StatusOK, *token)
}
