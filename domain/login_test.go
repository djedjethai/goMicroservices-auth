package domain

import (
	"database/sql"
	"fmt"
	// "reflect"
	"testing"
)

func Test_Login_return_claimForUser_if_Account_and_CustomerId_are_valid(t *testing.T) {
	// Arrange
	sqlCustID := sql.NullString{
		String: "10002",
		Valid:  true,
	}
	sqlAccID := sql.NullString{
		String: "1111",
		Valid:  true,
	}

	log := Login{
		CustomerId: sqlCustID,
		Accounts:   sqlAccID,
	}

	// Act
	claim := log.ClaimsForAccessToken()

	// Assert
	fmt.Printf("olalal: %v\n", claim.Accounts)
	if claim.CustomerId != "10002" && claim.Accounts[0] != "1111" {
		t.Error("while testing Login ClaimsForAccessToken should return claims for User if CustomerId and AccountId")
	}
}

func Test_Login_return_claimsForAdmin_if_CustomerId_is_not_valid(t *testing.T) {

	// Arrange
	sqlCustID := sql.NullString{
		String: "",
		Valid:  false,
	}
	sqlAccID := sql.NullString{
		String: "1111",
		Valid:  true,
	}

	log := Login{
		CustomerId: sqlCustID,
		Accounts:   sqlAccID,
	}

	// Act
	claim := log.ClaimsForAccessToken()

	// Assert
	fmt.Printf("olalal: %v\n", claim.Accounts)
	if len(claim.Accounts) > 0 {
		t.Error("while testing Login ClaimsForAccessToken should return claims for Admin if AccountId are not Valid")
	}

}

func Test_Login_return_claimsForAdmin_if_AccountID_is_not_valid(t *testing.T) {

	// Arrange
	sqlCustID := sql.NullString{
		String: "10002",
		Valid:  true,
	}
	sqlAccID := sql.NullString{
		String: "",
		Valid:  false,
	}

	log := Login{
		CustomerId: sqlCustID,
		Accounts:   sqlAccID,
	}

	// Act
	claim := log.ClaimsForAccessToken()

	// Assert
	if claim.CustomerId != "" {
		t.Error("while testing Login ClaimsForAccessToken should return claims for Admin if AccountId are not Valid")
	}

}
