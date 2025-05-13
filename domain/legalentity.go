package domain

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	//"gorm.io/datatypes"

	"github.com/google/uuid"
)

type JSONB map[string]interface{}

// Value Преобразует JSONB в значение для базы данных
func (j JSONB) Value() (driver.Value, error) {
	if j == nil {
		return "{}", nil
	}
	return json.Marshal(j)
}

// Scan Преобразует значение из базы данных в JSONB
func (j *JSONB) Scan(value interface{}) error {
	if value == nil {
		*j = JSONB{}
		return nil
	}

	var data []byte
	switch v := value.(type) {
	case []byte:
		data = v
	case string:
		data = []byte(v)
	default:
		return fmt.Errorf("неподдерживаемый тип: %T", value)
	}

	return json.Unmarshal(data, j)
}

type LegalEntity struct {
	UUID uuid.UUID `gorm:"type:uuid;primaryKey" json:"uuid"`
	Name string    `gorm:"type:varchar(255);not null" validate:"required,lte=100" json:"name"`

	//CreatedBy     *string           `gorm:"type:varchar(255)" json:"created_by,omitempty"`
	//CreatedByUUID *uuid.UUID        `gorm:"type:uuid" json:"created_by_uuid,omitempty"`
	//Meta          datatypes.JSONMap `gorm:"type:jsonb" json:"meta"`

	BankAccounts []BankAccount `gorm:"foreignKey:LegalEntityID" json:"bank_accounts,omitempty"`

	CreatedAt time.Time  `gorm:"not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time  `gorm:"not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

func NewLegalEntity(name string) *LegalEntity {
	entity := &LegalEntity{
		UUID: uuid.New(),
		Name: name,
		//Meta:      map[string]interface{}{},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
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
