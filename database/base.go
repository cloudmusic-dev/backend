package database

import (
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

type BaseModel struct {
	ID uuid.UUID `gorm:"type:uuid;primary_key;"`
}

func (base *BaseModel) BeforeCreate(scope *gorm.Scope) error {
	id, err := uuid.NewRandom()
	if err != nil {
		return err
	}
	return scope.SetColumn("ID", id)
}
