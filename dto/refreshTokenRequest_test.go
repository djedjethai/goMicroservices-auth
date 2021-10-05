package dto

import (
	"fmt"
	"testing"
)

func Test_dtoRefreshTokenRequest_isAccessTokenValidity_return_err_if_token_is_not_valid(t *testing.T) {
	// create token
	refreshToken := RefreshTokenRequest{
		AccessToken:  "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0b2tlbl90eXBlIjoicmVmcmVzaF90b2tlbiIsImN1c3RvbWVyX2lkIjoiIiwiYWNjb3VudHMiOm51bGwsInVuIjoiYWRtaW4iLCJyb2xlIjoiYWRtaW4iLCJleHAiOjE2MjkzODE4MjB9.dtsik_uSKfoduArFg0ZuneApz9IfNN0rOL1rS-ByuM8",
		RefreshToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0b2tlbl90eXBlIjoicmVmcmVzaF90b2tlbiIsImN1c3RvbWVyX2lkIjoiIiwiYWNjb3VudHMiOm51bGwsInVuIjoiYWRtaW4iLCJyb2xlIjoiYWRtaW4iLCJleHAiOjE2MjkzODE4MjB9.dtsik_uSKfoduArFg0ZuneApz9IfNN0rOL1rS-ByuM9",
	}

	// print("grrrr: ", validationErr)
	// call func
	validationErr := refreshToken.IsAccessTokenValid()

	fmt.Printf("grrrr111: %v\n", validationErr)

	// assert
	if validationErr != nil {
		t.Error("While testing dtoRefreshTokenRequest IsAccessTokenValid shoud not return an err if token is valid")
	}
}

// func Test_dtoRefreshTokenRequest_isAccessTokenValidity_return_nil_if_token_is_valid(t *testing.T){}
