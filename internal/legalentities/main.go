package legalentities

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/krisch/crm-backend/domain"
	"github.com/sirupsen/logrus"
)

type Service interface {
	GetAllLegalEntities() ([]domain.LegalEntity, error)
	CreateLegalEntity(entity *domain.LegalEntity) error
	UpdateLegalEntity(entity *domain.LegalEntity) error
	DeleteLegalEntity(id string) error

	GetLegalEntityByUUID(uuid uuid.UUID) (*domain.LegalEntity, error)

	GetAllBankAccounts(legalEntityUUID uuid.UUID) ([]domain.BankAccount, error)
	GetBankAccount(accountUUID uuid.UUID) (*domain.BankAccount, error)

	CreateBankAccount(account *domain.BankAccount) error
	UpdateBankAccount(account *domain.BankAccount) error
	DeleteBankAccount(accountUUID uuid.UUID) error
}

type ServiceImpl struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &ServiceImpl{repo: repo}
}

func (s *ServiceImpl) GetAllLegalEntities() ([]domain.LegalEntity, error) {
	return s.repo.GetAll()
}

func (s *ServiceImpl) CreateLegalEntity(entity *domain.LegalEntity) error {
	entity.UUID = uuid.New()
	entity.CreatedAt = time.Now()
	entity.UpdatedAt = time.Now()
	return s.repo.Create(entity)
}

func (s *ServiceImpl) UpdateLegalEntity(entity *domain.LegalEntity) error {
	return s.repo.Update(entity)
}

func (s *ServiceImpl) DeleteLegalEntity(id string) error {
	return s.repo.Delete(id)
}

func (s *ServiceImpl) GetAllBankAccounts(legalEntityUUID uuid.UUID) ([]domain.BankAccount, error) {
	return s.repo.GetAllBankAccounts(legalEntityUUID)
}

func (s *ServiceImpl) CreateBankAccount(account *domain.BankAccount) error {
	// Генерируем новый UUID, если он не задан
	if account.UUID == uuid.Nil {
		account.UUID = uuid.New()
	}

	// Устанавливаем даты создания/обновления
	account.CreatedAt = time.Now()
	account.UpdatedAt = time.Now()

	// Проверяем валидность данных
	if err := ValidateBankAccount(account); err != nil {
		return fmt.Errorf("validation error: %w", err)
	}

	// Если счет основной, сбрасываем флаг у других счетов
	if account.IsPrimary {
		if err := s.repo.ClearPrimaryFlag(account.LegalEntityID); err != nil {
			logrus.WithFields(logrus.Fields{
				"legal_entity_id": account.LegalEntityID,
				"error":           err,
			}).Error("Failed to clear primary flag")
			return err
		}
	}

	// Сохраняем в БД
	return s.repo.CreateBankAccount(account)
}

func (s *ServiceImpl) UpdateBankAccount(account *domain.BankAccount) error {
	// Валидация данных
	if err := ValidateBankAccount(account); err != nil {
		return fmt.Errorf("validation error: %w", err)
	}

	// Получаем текущий счет из БД для проверки
	existingAcc, err := s.repo.GetBankAccount(account.UUID)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"account_uuid": account.UUID,
			"error":        err,
		}).Error("Failed to get existing bank account")
		return fmt.Errorf("failed to get bank account: %w", err)
	}

	// Если счет стал основным, сбрасываем флаг у других счетов
	if account.IsPrimary && (!existingAcc.IsPrimary || existingAcc.LegalEntityID != account.LegalEntityID) {
		if err := s.repo.ClearPrimaryFlag(account.LegalEntityID); err != nil {
			logrus.WithFields(logrus.Fields{
				"legal_entity_id": account.LegalEntityID,
				"error":           err,
			}).Error("Failed to clear primary flag")
			return fmt.Errorf("failed to clear primary flag: %w", err)
		}
	}

	// Обновляем счет
	if err := s.repo.UpdateBankAccount(account); err != nil {
		logrus.WithFields(logrus.Fields{
			"account_uuid": account.UUID,
			"error":        err,
		}).Error("Failed to update bank account")
		return fmt.Errorf("failed to update bank account: %w", err)
	}

	return nil
}

func (s *ServiceImpl) DeleteBankAccount(accountUUID uuid.UUID) error {
	return s.repo.DeleteBankAccount(accountUUID)
}

func (s *ServiceImpl) GetBankAccount(accountUUID uuid.UUID) (*domain.BankAccount, error) {
	return s.repo.GetBankAccount(accountUUID)
}

func (s *ServiceImpl) GetLegalEntityByUUID(uuid uuid.UUID) (*domain.LegalEntity, error) {
	entities, err := s.repo.GetAll()
	if err != nil {
		return nil, fmt.Errorf("failed to get legal entities: %w", err)
	}

	for _, entity := range entities {
		if entity.UUID == uuid {
			return &entity, nil
		}
	}

	return nil, ErrLegalEntityNotFound
}
