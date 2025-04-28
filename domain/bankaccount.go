package domain

import (
	"log"
	"time"

	web "github.com/krisch/crm-backend/internal/web/olegalentities"

	"github.com/google/uuid"
)

type BankAccount struct {
	UUID          uuid.UUID `gorm:"primaryKey" json:"uuid"`
	LegalEntityID uuid.UUID `gorm:"type:uuid;not null" json:"legal_entity_id"`

	BIC                  string `gorm:"column:bic;not null" json:"bic"`
	BankName             string `gorm:"column:bank_name;not null" json:"bank_name"`
	BankAddress          string `gorm:"column:bank_address" json:"bank_address"`
	CorrespondentAccount string `gorm:"column:correspondent_account" json:"correspondent_account"`
	SettlementAccount    string `gorm:"column:settlement_account;not null" json:"settlement_account"` // Было CheckingAccount
	Currency             string `gorm:"default:'RUB'" json:"currency"`
	Comment              string `json:"comment"`
	IsPrimary            bool   `gorm:"default:false" json:"is_primary"`

	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func NewBankAccount(legalEntityID uuid.UUID, dto web.BankAccountDTO) *BankAccount {

	if dto.CheckingAccount == "" {
		log.Println("settlement_account is required")
		return nil
	}

	return &BankAccount{
		UUID:                 uuid.New(),
		LegalEntityID:        legalEntityID,
		BIC:                  dto.Bik,
		BankName:             dto.Bank,
		BankAddress:          getString(dto.Address),
		CorrespondentAccount: getString(dto.CorrespondentAccount),
		SettlementAccount:    dto.CheckingAccount,
		Currency:             "RUB",
		Comment:              getString(dto.Comment),
		IsPrimary:            getBool(dto.IsPrimary),
	}
}

func getString(s *string) string {
	if s != nil {
		return *s
	}
	return ""
}

func getBool(b *bool) bool {
	if b != nil {
		return *b
	}
	return false
}
