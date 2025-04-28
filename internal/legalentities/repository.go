package legalentities

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/krisch/crm-backend/domain"
	"gorm.io/gorm"
)

type Repository interface {
	GetAll() ([]domain.LegalEntity, error)
	Create(entity *domain.LegalEntity) error
	Update(entity *domain.LegalEntity) error
	Delete(uuid string) error

	GetAllBankAccounts(legalEntityUUID uuid.UUID) ([]domain.BankAccount, error)
	CreateBankAccount(account *domain.BankAccount) error
	UpdateBankAccount(account *domain.BankAccount) error
	DeleteBankAccount(accountUUID uuid.UUID) error
	GetBankAccount(accountUUID uuid.UUID) (*domain.BankAccount, error)
	ClearPrimaryFlag(legalEntityID uuid.UUID) error
}

func (r *repository) ClearPrimaryFlag(legalEntityID uuid.UUID) error {
	return r.db.Model(&domain.BankAccount{}).
		Where("legal_entity_id = ? AND is_primary = true", legalEntityID).
		Update("is_primary", false).
		Error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) GetAll() ([]domain.LegalEntity, error) {
	var entities []domain.LegalEntity
	if err := r.db.Preload("BankAccounts").Find(&entities).Error; err != nil {
		return nil, err
	}
	return entities, nil
}

func (r *repository) Create(entity *domain.LegalEntity) error {
	return r.db.Create(entity).Error
}

func (r *repository) Update(entity *domain.LegalEntity) error {
	return r.db.Save(entity).Error
}

func (r *repository) Delete(uuidStr string) error {
	parsedUUID, err := uuid.Parse(uuidStr)
	if err != nil {
		return err
	}
	return r.db.Where("uuid = ?", parsedUUID).Delete(&domain.LegalEntity{}).Error
}

func (r *repository) GetAllBankAccounts(legalEntityUUID uuid.UUID) ([]domain.BankAccount, error) {
	var accounts []domain.BankAccount
	result := r.db.Where("legal_entity_id = ?", legalEntityUUID).Find(&accounts)
	return accounts, result.Error
}

func (r *repository) CreateBankAccount(account *domain.BankAccount) error {
	// Генерируем UUID, если он не задан
	if account.UUID == uuid.Nil { // Проверка на "пустой" UUID
		account.UUID = uuid.New() // Генерация нового UUID
	}

	// Устанавливаем даты, если они не заданы
	if account.CreatedAt.IsZero() {
		account.CreatedAt = time.Now()
	}
	if account.UpdatedAt.IsZero() {
		account.UpdatedAt = time.Now()
	}

	// Если счет основной, сбрасываем флаг у других счетов
	if account.IsPrimary {
		if err := r.ClearPrimaryFlag(account.LegalEntityID); err != nil {
			return fmt.Errorf("failed to clear primary flag: %w", err)
		}
	}

	// Сохраняем через GORM
	if err := r.db.Create(account).Error; err != nil {
		return fmt.Errorf("failed to create bank account: %w", err)
	}

	return nil
}

func (r *repository) UpdateBankAccount(account *domain.BankAccount) error {
	if account.IsPrimary { // Было: if account.IsPrimary != nil && *account.IsPrimary
		if err := r.ClearPrimaryFlag(account.LegalEntityID); err != nil {
			return err
		}
	}
	return r.db.Save(account).Error
}

func (r *repository) DeleteBankAccount(accountUUID uuid.UUID) error {
	return r.db.Delete(&domain.BankAccount{}, "uuid = ?", accountUUID).Error
}

func (r *repository) GetBankAccount(accountUUID uuid.UUID) (*domain.BankAccount, error) {
	var account domain.BankAccount
	if err := r.db.Where("uuid = ?", accountUUID).First(&account).Error; err != nil {
		return nil, err
	}
	return &account, nil
}
