package dto

import (
	"github.com/djedjethai/bankingAuth/errs"
)

type SignupRequest struct {
	Username     string `json:"username"`
	Password     string `json:"password"`
	PasswordConf string `json:"password_conf"`
}

func (sr SignupRequest) ValidUsernameAndPwdSyntax() *errs.AppError {
	if len(sr.Username) < 5 || len(sr.Password) < 5 {
		return errs.NewValidationError("Username or Password Syntaxe is invalid")
	}

	if sr.Password != sr.PasswordConf {
		return errs.NewValidationError("Password and Password's confirmation do not match")
	}
}
