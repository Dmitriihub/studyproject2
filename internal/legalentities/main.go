package legalentities

import (
	"time"

	"github.com/google/uuid"
	"github.com/krisch/crm-backend/domain"
)

type Service interface {
	GetAllLegalEntities() ([]domain.LegalEntity, error)
	CreateLegalEntity(entity *domain.LegalEntity) error
	UpdateLegalEntity(entity *domain.LegalEntity) error
	DeleteLegalEntity(id string) error
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
	entity.UUID = uuid.New() // Преобразуем UUID в строку
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
