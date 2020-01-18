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
	ActivationCode string
	CreatedAt      time.Time

	Playlists 	   []Playlist `gorm:"foreignkey:Owner"`
}
