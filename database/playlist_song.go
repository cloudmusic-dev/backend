package database

import "github.com/google/uuid"

type PlaylistSong struct {
	PlaylistId uuid.UUID `gorm:"column:playlistId"`
	SongId     uuid.UUID `gorm:"column:songId"`
}
