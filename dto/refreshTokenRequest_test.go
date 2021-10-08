package dto

import (
	// "fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/djedjethai/bankingAuth/domain"
	"testing"
	"time"
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

	// fmt.Printf("grrrr111: %v\n", validationErr)

	// assert
	if validationErr == nil {
		t.Error("While testing dtoRefreshTokenRequest IsAccessTokenValid shoud return an err if token is invalid")
	}
}

func Test_dtoRefreshTokenRequest_isAccessTokenValidity_return_nil_if_token_is_valid(t *testing.T) {
	// Arrange
	// create a token
	rtc := domain.AccessTokenClaims{
		CustomerId: "2001",
		Username:   "2001",
		Role:       "user",
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(domain.REFRESH_TOKEN_DURATION).Unix(),
		},
	}

	// create a new and valid token
	at := domain.NewAuthToken(rtc)
	tkn, _ := at.Token.SignedString([]byte(domain.HMAC_SAMPLE_SECRET))

	rt := RefreshTokenRequest{
		AccessToken:  tkn,
		RefreshToken: tkn,
	}
	// run func with queried token
	errTok := rt.IsAccessTokenValid()

	// assert return from func is nil
	if errTok != nil {
		t.Error("While testing dtoRefreshTokenValue IsAccessTokenValid should not return an err if token is valid")
	}
}
