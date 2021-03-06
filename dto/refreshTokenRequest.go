package dto

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/djedjethai/bankingAuth/domain"
)

type RefreshTokenRequest struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (r RefreshTokenRequest) IsAccessTokenValid() *jwt.ValidationError {
	// 1.valid token
	// 2.valid token but expired
	_, err := jwt.Parse(r.AccessToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(domain.HMAC_SAMPLE_SECRET), nil
	})

	if err != nil {
		var vErr *jwt.ValidationError
		// it s a type assertion, it check if the err if of type ValidationError
		// we only return the err in case it's a ValidationError
		// means the token come from us but is expire
		// next step of the err handling into the service
		if errors.As(err, &vErr) {
			return vErr
		}
	}

	return nil

}
