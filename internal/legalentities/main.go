package legalentities

import "github.com/google/uuid"

type Service interface {
	GetAllLegalEntities() ([]LegalEntity, error)
	CreateLegalEntity(entity *LegalEntity) error
	UpdateLegalEntity(entity *LegalEntity) error
	DeleteLegalEntity(uuid string) error
}

type ServiceImpl struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &ServiceImpl{repo: repo}
}

func (s *ServiceImpl) GetAllLegalEntities() ([]LegalEntity, error) {
	return s.repo.GetAll()
}

func (s *ServiceImpl) CreateLegalEntity(entity *LegalEntity) error {
	entity.UUID = uuid.New().String()
	return s.repo.Create(entity)
}

func (s *ServiceImpl) UpdateLegalEntity(entity *LegalEntity) error {
	return s.repo.Update(entity)
}

func (s *ServiceImpl) DeleteLegalEntity(uuid string) error {
	return s.repo.Delete(uuid)
}
