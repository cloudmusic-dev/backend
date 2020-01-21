DROP TABLE playlist_songs;
ALTER TABLE songs ADD playlist varchar(36);
ALTER TABLE songs ADD FOREIGN KEY (playlist) REFERENCES playlists(id);