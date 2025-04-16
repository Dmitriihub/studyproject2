package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/krisch/crm-backend/internal/helpers"
)

type BankAccount struct {
	UUID          uuid.UUID `gorm:"primaryKey"`
	LegalEntityID uuid.UUID `gorm:"type:uuid"`
	AccountNumber string    `validate:"required"`
	BankName      string    `validate:"required"`
	Correspondent string
	BIC           string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func NewBankAccount(legalEntityID uuid.UUID, accountNumber, bankName string) *BankAccount {
	ba := &BankAccount{
		UUID:          uuid.New(),
		LegalEntityID: legalEntityID,
		AccountNumber: accountNumber,
		BankName:      bankName,
	}

	errs, ok := helpers.ValidationStruct(ba)
	if !ok {
		panic(errors.New(helpers.Join(errs, ", ")))
	}

	return ba
}
