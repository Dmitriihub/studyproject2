package legalentities

import "errors"

var (
	ErrBankAccountNotFound    = errors.New("bank account not found")
	ErrInvalidBankAccountData = errors.New("invalid bank account data")
	ErrPrimaryAccountExists   = errors.New("primary account already exists")
	ErrLegalEntityNotFound    = errors.New("legal entity not found")

	ErrInvalidBIC           = errors.New("invalid BIK format")
	ErrInvalidAccountNumber = errors.New("invalid account number format")
	ErrInvalidCurrency      = errors.New("invalid currency format")
)

type ValidationError struct {
	Field   string
	Message string
}

func (e ValidationError) Error() string {
	return e.Field + ": " + e.Message
}
