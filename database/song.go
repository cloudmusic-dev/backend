package database

import "github.com/google/uuid"

type Song struct {
	BaseModel
	Owner      uuid.UUID
	Title      string
	Duration   int
	AudioFile  *uuid.UUID `gorm:"column:audioFile"`
	CoverFiler *uuid.UUID `gorm:"column:coverFile"`
}
