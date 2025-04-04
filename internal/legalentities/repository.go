package legalentities

import "gorm.io/gorm"

type Repository interface {
	GetAll() ([]LegalEntity, error)
	Create(entity *LegalEntity) error
	Update(entity *LegalEntity) error
	Delete(uuid string) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) GetAll() ([]LegalEntity, error) {
	var entities []LegalEntity
	if err := r.db.Find(&entities).Error; err != nil {
		return nil, err
	}
	return entities, nil
}

func (r *repository) Create(entity *LegalEntity) error {
	return r.db.Create(entity).Error
}

func (r *repository) Update(entity *LegalEntity) error {
	return r.db.Save(entity).Error
}

func (r *repository) Delete(uuid string) error {
	return r.db.Delete(&LegalEntity{UUID: uuid}).Error
}
