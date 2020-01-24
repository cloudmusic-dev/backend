package database

import (
	"database/sql"
	"time"
)

type User struct {
	BaseModel
	Username       string
	Email          string
	Password       string
	Activated      sql.NullBool
	ActivationCode string `gorm:"column:activationCode"`
	CreatedAt      time.Time `gorm:"column:createdAt"`

	Playlists []Playlist `gorm:"foreignkey:Owner"`
}
