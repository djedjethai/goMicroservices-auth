package dto

import (
	"github.com/djedjethai/bankingAuth/errs"
	"time"
)

type SignupRequest struct {
	Name         string    `json:"name"`
	DateOfBirth  time.Time `json:"date_of_birth"`
	City         string    `json:"city"`
	ZipCode      int       `json:"zip_code"`
	Username     string    `json:"username"`
	Password     string    `json:"password"`
	PasswordConf string    `json:"password_conf"`
}

func (sr SignupRequest) ValidNameDobCityZip() *errs.AppError {
	if len(sr.Name) < 4 ||
		len(sr.City < 2) ||
		len(sr.ZipCode) < 5 ||
		len(sr.ZipCode) > 5 {
		return errs.NewValidationError("Name or City or ZipCode is invalid")
	}

	if sr.DateOfBirth.IsZero {
		return errs.NewValidationError("Date of bith incorrect")
	}
	return nil
}

func (sr SignupRequest) ValidNewUser() *errs.AppError {
	if len(sr.Username) < 5 || len(sr.Password) < 5 {
		return errs.NewValidationError("Username or Password Syntaxe is invalid")
	}

	if sr.Password != sr.PasswordConf {
		return errs.NewValidationError("Password and Password's confirmation do not match")
	}
	return nil
}
