package api

import "github.com/cloudmusic-dev/backend/database"

type Song struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Duration int    `json:"duration"`
}

func databaseSongToApi(song database.Song) Song {
	ret := Song{
		ID: song.ID.String(),
		Title: song.Title,
		Duration: song.Duration,
	}

	return ret
}