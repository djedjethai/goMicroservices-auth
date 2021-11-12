package domain

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/djedjethai/bankingLib/errs"
	"github.com/djedjethai/bankingLib/logger"
)

// // to run the mock: go generate ./...
//
// //go:generate mockgen -destination=../mocks/domain/mockAuthToken.go -package=domain github.com/djedjethai/bankingAuth/domain AuthTokenInterface
// type AuthTokenInterface interface {
// 	NewAccessTokenFromRefreshToken(string) (string, *errs.AppError)
// 	NewAccessToken() (string, *errs.AppError)
// 	NewRefreshToken() (string, *errs.AppError)
// }

type AuthToken struct {
	Token *jwt.Token
}

func NewAuthToken(claims AccessTokenClaims) AuthToken {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return AuthToken{Token: token}
}

func NewAccessTokenFromRefreshToken(refreshToken string) (string, *errs.AppError) {
	token, err := jwt.ParseWithClaims(refreshToken, &RefreshTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(HMAC_SAMPLE_SECRET), nil
	})
	if err != nil {
		return "", errs.NewAuthenticationError("invalid or expired refresh token")
	}
	// get the claim, and we cast it as we know it's this type
	r := token.Claims.(*RefreshTokenClaims)
	accessTokenClaims := r.AccessTokenClaims()
	authToken := NewAuthToken(accessTokenClaims)

	return authToken.NewAccessToken()
}

func (t AuthToken) NewAccessToken() (string, *errs.AppError) {
	signedString, err := t.Token.SignedString([]byte(HMAC_SAMPLE_SECRET))
	if err != nil {
		logger.Error("Failed while signing access token: " + err.Error())
		return "", errs.NewInternalServerError("cannot generate access token")
	}
	return signedString, nil
}

func (t AuthToken) NewRefreshToken() (string, *errs.AppError) {
	// get the claim which is already store in the t.Token
	// from when we generated the accessToken
	// (as this token is generated immediately the accessToken, to be return in the same resp)
	c := t.Token.Claims.(AccessTokenClaims)
	refreshClaims := c.RefreshTokenClaims()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	signedString, err := token.SignedString([]byte(HMAC_SAMPLE_SECRET))
	if err != nil {
		logger.Error("Unable to generate a refresh token")
		return "", errs.NewInternalServerError("An Unexpected error occured")
	}

	return signedString, nil
}
