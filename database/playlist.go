package database

import (
	"github.com/google/uuid"
	"time"
)

type Playlist struct {
	BaseModel
	Owner        uuid.UUID
	Name         string
	ParentFolder uuid.UUID `gorm:"column:parentFolder"`
	CreatedAt    time.Time `gorm:"column:createdAt"`
}
