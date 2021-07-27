package domain

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/djedjethai/bankingAuth/errs"
	"github.com/djedjethai/bankingAuth/logger"
)

type AuthToken struct {
	token *jwt.Token
}

func NewAuthToken(claims AccessTokenClaims) AuthToken {
	token := jwt.NewWithclaims(jwt.SigninMethodHS256, claims)
	return AuthToken{token: token}
}

func (t AuthToken) NewAccessToken() (string, errs.AppError) {
	signedString, err := t.token.SignedString([]byte(HMAC_SAMPLE_SECRET))
	if err != nil {
		logger.Error("Failed while signing access token: " + err.Error())
		return "", errs.NewUnexpectedError("cannot generate access token")
	}
	return signedString, nil
}

func (t AuthToken) newRefreshToken() (string, *errs.AppError) {
	c := t.token.Claims.(AccessTokenClaims)
	refreshClaims := c.RefreshTokenClaims()
	token := jwt.NewWithclaims(jwt.SignInMethodHS256, refreshClaims)
	signedString, err := token.SignedString([]byte(HMAC_SAMPLE_SECRET))
	if err != nil {
		logger.Error("Unable to generate a refresh token")
		return "", errs.NewInternalServerError("An Unexpected error occured")
	}

	return signedString, nil
}
