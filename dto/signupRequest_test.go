package dto

import (
	// "fmt"
	"testing"
)

func Test_if_dtoSignupRequest_ValidNameDobCityZip_wrong_name_return_an_err(t *testing.T) {
	// Arrange
	sr := SignupRequest{
		Name:        "jer",
		City:        "bangkok",
		DateOfBirth: "1970-01-02",
		ZipCode:     "10000",
	}

	// Act
	err := sr.ValidNameDobCityZip()

	// Assert
	if err == nil {
		t.Error("While testing dtoSignupRequest ValidNameDobCityZip err should not be null")
	}
}

func Test_if_dtoSignupRequest_ValidNameDobCityZip_wrong_City_return_an_err(t *testing.T) {
	// Arrange
	sr := SignupRequest{
		Name:        "jerome",
		City:        "b",
		DateOfBirth: "1970-01-02",
		ZipCode:     "10000",
	}

	// Act
	err := sr.ValidNameDobCityZip()

	// Assert
	if err == nil {
		t.Error("While testing dtoSignupRequest ValidNameDobCityZip err should not be null")
	}
}

func Test_if_dtoSignupRequest_ValidNameDobCityZip_wrong_Zip_return_an_err(t *testing.T) {
	// test zip code too short
	sr := SignupRequest{
		Name:        "jerome",
		City:        "bangkok",
		DateOfBirth: "1970-01-02",
		ZipCode:     "1000",
	}
	err := sr.ValidNameDobCityZip()
	if err == nil {
		t.Error("While testing dtoSignupRequest ValidNameDobCityZip err should not be null")
	}

	// test zip code too long
	sr = SignupRequest{
		Name:        "jerome",
		City:        "bangkok",
		DateOfBirth: "1970-01-02",
		ZipCode:     "100000",
	}
	err = sr.ValidNameDobCityZip()
	if err == nil {
		t.Error("While testing dtoSignupRequest ValidNameDobCityZip err should not be null")
	}
}

func Test_if_dtoSignupRequest_ValidNameDobCityZip_return_nil(t *testing.T) {
	// Arrange
	sr := SignupRequest{
		Name:        "jerome",
		City:        "bankk",
		DateOfBirth: "1970-01-02",
		ZipCode:     "10000",
	}

	// Act
	err := sr.ValidNameDobCityZip()

	// Assert
	if err != nil {
		t.Error("While testing dtoSignupRequest ValidNameDobCityZip err be null")
	}
}

func Test_if_dtoSignupRequest_ValidNewUser_wrong_Username_return_an_err(t *testing.T) {
	// test
	sr := SignupRequest{
		Name:         "jerome",
		City:         "bankk",
		DateOfBirth:  "1970-01-02",
		ZipCode:      "10000",
		Username:     "jero",
		Password:     "password",
		PasswordConf: "password",
	}
	err := sr.ValidNewUser()
	if err == nil {
		t.Error("While testing dtoSignupRequest ValidNewUser err not be null")
	}
}

func Test_if_dtoSignupRequest_ValidNewUser_wrong_Password_return_an_err(t *testing.T) {
	sr := SignupRequest{
		Name:         "jerome",
		City:         "bankk",
		DateOfBirth:  "1970-01-02",
		ZipCode:      "10000",
		Username:     "eromme",
		Password:     "pass",
		PasswordConf: "pass",
	}
	err := sr.ValidNewUser()
	if err == nil {
		t.Error("While testing dtoSignupRequest ValidNewUser err not be null")
	}
}

func Test_if_dtoSignupRequest_ValidNewUser_password_and_not_match_passwordConf_return_an_err(t *testing.T) {
	sr := SignupRequest{
		Name:         "jerome",
		City:         "bankk",
		DateOfBirth:  "1970-01-02",
		ZipCode:      "10000",
		Username:     "eromme",
		Password:     "password",
		PasswordConf: "passwordnotmatch",
	}
	err := sr.ValidNewUser()
	if err == nil {
		t.Error("While testing dtoSignupRequest ValidNewUser err not be null")
	}
}

func Test_if_dtoSignupRequest_ValidNewUser_return_nil(t *testing.T) {
	sr := SignupRequest{
		Name:         "jerome",
		City:         "bankk",
		DateOfBirth:  "1970-01-02",
		ZipCode:      "10000",
		Username:     "eromme",
		Password:     "password",
		PasswordConf: "password",
	}
	err := sr.ValidNewUser()
	if err != nil {
		t.Error("While testing dtoSignupRequest ValidNewUser err should be null")
	}
}
