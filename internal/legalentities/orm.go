package legalentities

import (
	"github.com/krisch/crm-backend/domain"
	"gorm.io/gorm"
)

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&domain.LegalEntity{},
		&domain.BankAccount{},
	)
}
