package legalentities

import (
	"github.com/google/uuid"
	"github.com/krisch/crm-backend/domain"
	"gorm.io/gorm"
)

type Repository interface {
	GetAll() ([]domain.LegalEntity, error)
	Create(entity *domain.LegalEntity) error
	Update(entity *domain.LegalEntity) error
	Delete(uuid string) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) GetAll() ([]domain.LegalEntity, error) {
	var entities []domain.LegalEntity
	if err := r.db.Find(&entities).Error; err != nil {
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
	return r.db.Delete(&domain.LegalEntity{UUID: parsedUUID}).Error
}
