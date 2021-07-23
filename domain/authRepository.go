package domain

import (
	"database/sql"
	"github.com/djedjethai/bankingAuth/errs"
	"github.com/jmoiron/sqlx"
)

type AuthRepository interface {
}

type authRepository struct {
	client *sqlx.DB
}

func NewAuthRepository(client *sqlx.DB) AuthRepository {
	return authRepository{client}
}

func (c authRepository) FindBy(username, password string) (*Login, *errs.AppError) {
	var login Login
	sqlVerify := `SELECT username, u.customer_id, role, group_concat(a.account_id) as account_numbers FROM users u 
	LEFT JOIN account a ON a.customer_id = u.customer_id
	WHERE username = ? AND password = ?
	GROUP BY a.customer_id`
	err := c.client.Get(&login, sqlVerify, username, password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errs.NewValidationError("invalid credentials")
		} else {
			return nil, errs.NewInternalServerError("Unexpected database error")
		}
	}

	return &login, nil
}
