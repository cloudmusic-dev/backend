package database

import "time"

type User struct {
	BaseModel
	Username       string
	Email          string
	Password       string
	Activated      bool
	ActivationCode string
	CreatedAt      time.Time

	Playlists 	   []Playlist `gorm:"foreignkey:Owner"`
}
