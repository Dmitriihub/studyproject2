package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/krisch/crm-backend/internal/helpers"
	"gorm.io/datatypes"
)

type LegalEntity struct {
	UUID uuid.UUID `gorm:"type:uuid;primaryKey" json:"uuid"`
	Name string    `gorm:"type:varchar(255);not null" validate:"required,lte=100" json:"name"`

	CreatedBy     *string        `gorm:"type:varchar(255)" validate:"omitempty,email" json:"created_by,omitempty"`
	CreatedByUUID *uuid.UUID     `gorm:"type:uuid" validate:"omitempty,uuid" json:"created_by_uuid,omitempty"`
	Meta          datatypes.JSON `gorm:"type:jsonb" json:"meta"`

	CreatedAt time.Time  `gorm:"not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time  `gorm:"not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

func NewLegalEntity(name string, createdBy *string, createdByUUID *uuid.UUID) *LegalEntity {
	entity := &LegalEntity{ // Добавляем присваивание переменной
		UUID:          uuid.New(),
		Name:          name,
		CreatedBy:     createdBy,
		CreatedByUUID: createdByUUID,
		Meta:          datatypes.JSON([]byte("{}")),
	}

	errs, ok := helpers.ValidationStruct(entity)
	if !ok {
		panic(errors.New(helpers.Join(errs, ", ")))
	}

	return entity
}

func NewLegalEntityUUID(uid uuid.UUID) *LegalEntity {
	return &LegalEntity{
		UUID: uid,
	}
}

func (e *LegalEntity) ChangeName(name string) error {
	if len(name) < 1 || len(name) > 100 {
		return errors.New("название должно быть от 1 до 100 символов")
	}
	e.Name = name
	return nil
}
