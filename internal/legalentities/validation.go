package legalentities

import (
	"fmt"
	"regexp"

	"github.com/krisch/crm-backend/domain"
)

var (
	bikRegexp      = regexp.MustCompile(`^\d{9}$`)
	accountRegexp  = regexp.MustCompile(`^\d{20}$`)
	corrAccRegexp  = regexp.MustCompile(`^\d{20}$`)
	currencyRegexp = regexp.MustCompile(`^[A-Z]{3}$`)
)

func ValidateBankAccount(acc *domain.BankAccount) error {
	if !bikRegexp.MatchString(acc.BIC) {
		return fmt.Errorf("invalid BIK format")
	}

	if !accountRegexp.MatchString(acc.SettlementAccount) {
		return fmt.Errorf("invalid settlement account format")
	}

	if acc.CorrespondentAccount != "" && !corrAccRegexp.MatchString(acc.CorrespondentAccount) {
		return fmt.Errorf("invalid correspondent account format")
	}

	if acc.Currency != "" && !currencyRegexp.MatchString(acc.Currency) {
		return fmt.Errorf("invalid currency format")
	}

	if len(acc.BankName) > 255 {
		return fmt.Errorf("bank name too long")
	}

	return nil
}
