CREATE TABLE playlist_songs(
    playlistId varchar(36) not null,
    songId varchar(36) not null,
    primary key (playlistId, songId),
    foreign key (playlistId) references playlists(id),
    foreign key (songId) references songs(id)
);
ALTER TABLE songs DROP FOREIGN KEY songs_ibfk_2;
ALTER TABLE songs DROP COLUMN playlist;