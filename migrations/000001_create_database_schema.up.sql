CREATE TABLE users
(
	id varchar(36) not null primary key,
	username varchar(50) not null,
	email varchar(255) not null,
	password char(60) not null,
	activated bit default 0 not null,
	activationCode varchar(100) not null,
	createdAt datetime default now(),
	unique key unique_username(username),
	unique key unique_email(email)
);

CREATE TABLE folders
(
    id varchar(36) not null primary key,
    owner varchar(36) not null,
    name varchar(150) not null,
    parentFolder varchar(36) null,
    createdAt datetime default now(),
    foreign key (owner) references users(id),
    foreign key (parentFolder) references folders(id)
);

CREATE TABLE playlists
(
    id varchar(36) not null primary key,
    owner varchar(36) not null,
    name varchar(150) not null,
    parentFolder varchar(36) null,
    createdAt datetime default now(),
    foreign key (owner) references users(id),
    foreign key (parentFolder) references folders(id)
);

CREATE TABLE storage_items
(
    id varchar(36) not null primary key,
    storage varchar(36) not null,
    path varchar(255) not null,
    originalName varchar(255) not null
);

CREATE TABLE songs
(
    id varchar(36) not null primary key,
    owner varchar(36) not null,
    title varchar(255) not null,
    duration int not null,
    playlist varchar(36) null,
    audioFile varchar(36) null,
    coverFile varchar(36) null,
    foreign key (owner) references users(id),
    foreign key (playlist) references playlists(id),
    foreign key (audioFile) references storage_items(id),
    foreign key (coverFile) references storage_items(id)
);

CREATE TABLE song_artists
(
    song varchar(36) not null,
    artist varchar(255) not null,
    primary key (song, artist)
);