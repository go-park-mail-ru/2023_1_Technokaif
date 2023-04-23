CREATE TYPE SEX AS ENUM('M', 'F', 'O'); 
CREATE TABLE Users
(
    id            SERIAL       PRIMARY KEY,
    version       INT                      NOT NULL DEFAULT 1,
    username      VARCHAR(20)   UNIQUE     NOT NULL,
    email         VARCHAR(255)  UNIQUE     NOT NULL,
    password_hash VARCHAR(256)             NOT NULL,
    salt          VARCHAR(64)              NOT NULL,
    first_name    VARCHAR(20)              NOT NULL,
    last_name     VARCHAR(20)              NOT NULL,
    sex           SEX                      NOT NULL,
    birth_date    DATE                     NOT NULL,
    avatar_src    TEXT
);

CREATE TABLE Artists
(
    id         SERIAL      PRIMARY KEY,
    user_id    INT         REFERENCES Users(id) ON DELETE SET NULL,
    name       VARCHAR(30)                                         NOT NULL,
    avatar_src TEXT                                                NOT NULL
);

CREATE TABLE Albums
(
    id          SERIAL        PRIMARY KEY,
    name        VARCHAR(40)               NOT NULL,
    description VARCHAR(2000),
    cover_src   TEXT                      NOT NULL
);

CREATE TABLE Artists_Albums
(
    artist_id INT REFERENCES Artists(id) ON DELETE CASCADE NOT NULL,
    album_id  INT REFERENCES Albums(id)  ON DELETE CASCADE NOT NULL,

    PRIMARY KEY(artist_id, album_id)
);

CREATE TABLE Tracks
(
    id             SERIAL      PRIMARY KEY,
    name           VARCHAR(40)                                          NOT NULL,
    album_id       INT         REFERENCES Albums(id)  ON DELETE CASCADE,
    album_position INT,
    cover_src      TEXT                                                 NOT NULL,
    record_src     TEXT                                                 NOT NULL,
    listens        INT         DEFAULT 0                                NOT NULL,

    UNIQUE(album_id, album_position)
);

CREATE TABLE Listens
(
    id          SERIAL      PRIMARY KEY,
    user_id     INT         REFERENCES  Users(id)  ON DELETE SET NULL,
    track_id    INT         REFERENCES  Tracks(id) ON DELETE CASCADE  NOT NULL,
    commited_at TIMESTAMPTZ DEFAULT NOW()                             NOT NULL
);

CREATE TABLE Artists_Tracks
(
    artist_id INT REFERENCES Artists(id) ON DELETE CASCADE NOT NULL,
    track_id  INT REFERENCES Tracks(id)  ON DELETE CASCADE NOT NULL,

    PRIMARY KEY(artist_id, track_id)
);

CREATE TABLE Playlists
(
    id          SERIAL        PRIMARY KEY,
    name        VARCHAR(60)               NOT NULL,
    description VARCHAR(2000),
    cover_src   TEXT
);

CREATE TABLE Users_Playlists
(
    user_id     INT REFERENCES Users(id)     ON DELETE SET NULL,
    playlist_id INT REFERENCES Playlists(id) ON DELETE CASCADE,

    PRIMARY KEY(user_id, playlist_id)
);

CREATE TABLE Playlists_Tracks
(
    playlist_id INT REFERENCES Playlists(id) ON DELETE CASCADE,
    track_id    INT REFERENCES Tracks        ON DELETE CASCADE,
    added_at    TIMESTAMPTZ DEFAULT NOW()                      NOT NULL,

    PRIMARY KEY(playlist_id, track_id)
);

CREATE TABLE Liked_albums
(
    user_id   INT REFERENCES Users(id)  ON DELETE CASCADE NOT NULL,
    album_id  INT REFERENCES Albums(id) ON DELETE CASCADE NOT NULL,

    PRIMARY KEY(user_id, album_id)
);

CREATE TABLE Liked_artists
(
    user_id    INT REFERENCES Users(id)   ON DELETE CASCADE NOT NULL,
    artist_id  INT REFERENCES Artists(id) ON DELETE CASCADE NOT NULL,

    PRIMARY KEY(user_id, artist_id)
);

CREATE TABLE Liked_tracks
(
    user_id   INT REFERENCES Users(id)  ON DELETE CASCADE NOT NULL,
    track_id  INT REFERENCES Tracks(id) ON DELETE CASCADE NOT NULL,

    PRIMARY KEY(user_id, track_id)
);

CREATE TABLE Liked_playlists
(
    user_id     INT REFERENCES Users(id)     ON DELETE CASCADE NOT NULL,
    playlist_id INT REFERENCES Playlists(id) ON DELETE CASCADE NOT NULL,

    PRIMARY KEY(user_id, playlist_id)
);
