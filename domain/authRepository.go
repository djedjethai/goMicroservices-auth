package domain

import (
	"database/sql"
	// "fmt"
	"github.com/djedjethai/bankingAuth/errs"
	"github.com/djedjethai/bankingAuth/logger"
	"github.com/jmoiron/sqlx"
)

type AuthRepository interface {
	FindBy(string, string) (*Login, *errs.AppError)
	GenerateAndSaveRefreshTokenToStore(AuthToken) (string, *errs.AppError)
	RefreshTokenExists(string) *errs.AppError
	IsUsernameExist(string) (bool, *errs.AppError)
}

type authRepository struct {
	client *sqlx.DB
}

func NewAuthRepository(client *sqlx.DB) AuthRepository {
	return authRepository{client}
}

func (c authRepository) RefreshTokenExists(refreshToken string) *errs.AppError {
	logger.Info("the refresh token before the db Refresh: " + refreshToken)

	sqlSelect := "select refresh_token from refresh_token_store where refresh_token = ?"
	var token string
	err := c.client.Get(&token, sqlSelect, refreshToken)
	if err != nil {
		if err == sql.ErrNoRows {
			return errs.NewAuthenticationError("refresh token not registered in the store")
		} else {
			logger.Error("Unexpected database error: " + err.Error())
			return errs.NewInternalServerError("Unexpected database error")
		}
	}
	return nil
}

// here not goood need add to customer table first then user table:wq
func (c authRepository) CreateCustAndUser(cust CustomerDomain) (*Login, *errs.AppError) {

	var login Login

	tx, err := c.client.Begin()
	if err != nil {
		return nil, NewInternalServerError("Unexpected database error")
	}

	// insert into customer table first

	// insert into user table(using the id from customer table)
	result, errEx := tx.Exec(`INSERT INTO users (username, password, role, customer_id) values (?,?,?,?)`, username, password, "user", "--CustomerId--")

	// c.FindBy(username, password string) and get back the token
	// return the token

	// make sure it rollback all queries in case of one fail
	if errEx != nil {
		tx.Rollback()
		logger.Error("Error while craeting a new user" + errEx.Error())
		return nil, NewInternalServerError("Unexpected database error")
	}

}

func (c authRepository) IsUsernameExist(username string) (bool, *errs.AppError) {
	var name string

	sqlVerify := `SELECT username FROM users WHERE username = ?`
	err := c.client.Get(&name, sqlVerify, username)
	if err != nil {
		if err == sql.ErrNoRows {
			return true, nil
		} else {
			return false, errs.NewInternalServerError("Unexpected database error")
		}
	}

	return false, nil
}

func (c authRepository) FindBy(username, password string) (*Login, *errs.AppError) {
	var login Login

	sqlVerify := `SELECT username, u.customer_id, role, group_concat(a.account_id) as account_numbers FROM users u
                LEFT JOIN accounts a ON a.customer_id = u.customer_id
                WHERE username = ? and password = ?
                GROUP BY a.customer_id`
	err := c.client.Get(&login, sqlVerify, username, password)
	if err != nil {
		if err == sql.ErrNoRows {
			logger.Error("ErrNoRows in domain.FindBy " + err.Error())
			return nil, errs.NewValidationError("invalid credentials")
		} else {
			logger.Error("Err in domain.FindBy" + err.Error())
			return nil, errs.NewInternalServerError("Unexpected database error")
		}
	}

	return &login, nil
}

func (c authRepository) GenerateAndSaveRefreshTokenToStore(authToken AuthToken) (string, *errs.AppError) {
	// generate the refresh token
	var appErr *errs.AppError
	var refreshedToken string
	if refreshedToken, appErr = authToken.newRefreshToken(); appErr != nil {
		return "", appErr
	}

	// store it the store
	sqlInsert := "insert into refresh_token_store (refresh_token) values (?)"
	_, err := c.client.Exec(sqlInsert, refreshedToken)
	if err != nil {
		logger.Error("Unexpected database error when saving refresh token" + err.Error())
		return "", errs.NewInternalServerError("Unexpected server error")
	}

	return refreshedToken, nil
}
