package domain

import (
	"testing"
)

func Test_claim_IsUserRole_return_true_if_role_user_and_false_if_not(t *testing.T) {
	// Arrange
	roleAdmin := AccessTokenClaims{Role: "admin"}
	roleUser := AccessTokenClaims{Role: "user"}

	// Act
	adm := roleAdmin.IsUserRole()
	usr := roleUser.IsUserRole()

	//  Assert
	if adm != false {
		t.Error("While testing Claim IsUserRole should return false if it is not user")
	}

	if usr != true {
		t.Error("While testing Claim IsUserRole should return true if it is user")
	}
}

func Test_claim_IsValidCustomerId_return_true_customerId_is_valid_and_false_if_not(t *testing.T) {
	// Assert
	goodId := AccessTokenClaims{CustomerId: "2001"}
	badId := AccessTokenClaims{CustomerId: "200155"}

	// Act
	gi := goodId.IsValidCustomerId("2001")
	bi := badId.IsValidCustomerId("2001")

	// Assert
	if gi != true {
		t.Error("While testing Claim IsValidCustomerId valid CustometId should return true")
	}

	if bi != false {
		t.Error("While testing Claim IsValidCustomerId invalid CustometId should return false")
	}

}

func Test_claim_IsValidAccountId_return_false_if_account_not_found_and_true_if_it_is(t *testing.T) {
	// Assert
	act := AccessTokenClaims{Accounts: []string{"1111", "2222"}}

	// Act
	goodActId := act.IsValidAccountId("2222")
	badActId := act.IsValidAccountId("3333")

	// Assert
	if goodActId != true {
		t.Error("While testing Claim IsValidAccountId valid accountId should return true")
	}

	if badActId != false {
		t.Error("While testing Claim IsValidAccountId invalid accountId should return false")
	}

}

// A FINIR
func Test_claim_IsRequestVerifiedWithTokenClaims_return_false_CustumerId_or_IsValidAccountId_return_false_or_true_if_all_match(t *testing.T) {
	// Arrange
	cid := AccessTokenClaims{CustomerId: "10002", Accounts: []string{"0000"}}
	aid := AccessTokenClaims{CustomerId: "10000", Accounts: []string{"1111", "2222"}}
	caid := AccessTokenClaims{CustomerId: "10000", Accounts: []string{"0000"}}

	urlPar := make(map[string]string)
	urlPar["customer_id"] = "10000"
	urlPar["account_id"] = "0000"

	// Act
	custId := cid.IsRequestVerifiedWithTokenClaims(urlPar)
	accId := aid.IsRequestVerifiedWithTokenClaims(urlPar)
	custAccId := caid.IsRequestVerifiedWithTokenClaims(urlPar)

	// Assert
	if custId != false {
		t.Error("When testing Claim IsRequestVerifiedWithTokenClaims should return false if CustomerId does not match")
	}

	if accId != false {
		t.Error("When testing Claim IsRequestVerifiedWithTokenClaims should return false if AccountId does not match")
	}

	if custAccId != true {
		t.Error("When testing Claim IsRequestVerifiedWithTokenClaims should return true if AccountId and CustomerId match")
	}

}
