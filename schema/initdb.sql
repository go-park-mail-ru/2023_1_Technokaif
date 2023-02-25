CREATE TYPE SEX AS ENUM('M', 'F', 'O'); 
CREATE TABLE Users
(
    id            SERIAL       PRIMARY KEY,
    username      VARCHAR(20)  UNIQUE      NOT NULL,
    email         VARCHAR(30)  UNIQUE      NOT NULL,
    password_hash VARCHAR(256)             NOT NULL,
    first_name    VARCHAR(20)              NOT NULL,
    last_name     VARCHAR(20)              NOT NULL,
    user_sex      SEX                      NOT NULL,
    birth_date    DATE                     NOT NULL,
    avatar_src    TEXT
);

CREATE TABLE Artists
(
    id         SERIAL      PRIMARY KEY,
    name       VARCHAR(30)             NOT NULL,
    avatar_src TEXT
);

CREATE TABLE Albums
(
    id        SERIAL        PRIMARY KEY,
    name      VARCHAR(40)               NOT NULL,
    cover_src TEXT,
    info      VARCHAR(2000)
);

CREATE TABLE Tracks
(
    id         SERIAL      PRIMARY KEY,
    name       VARCHAR(40)                                          NOT NULL,
    album_id   INT         REFERENCES Albums(id)  ON DELETE CASCADE,
    artist_id  INT         REFERENCES Artists(id) ON DELETE CASCADE NOT NULL,
    cover_src  TEXT,
    record_src TEXT
);

CREATE TABLE Listens
(
    id       SERIAL PRIMARY KEY,
    user_id  INT    REFERENCES Users(id)  ON DELETE SET NULL,
    track_id INT    REFERENCES Tracks(id) ON DELETE CASCADE  NOT NULL
);

CREATE TABLE ArtistsAlbums
(   
    artist_id INT REFERENCES Artists(id) ON DELETE CASCADE NOT NULL,
    album_id  INT REFERENCES Albums(id)  ON DELETE CASCADE NOT NULL,

    PRIMARY KEY(artist_id, album_id)
);

CREATE TABLE ArtistsTracks
(
    artist_id INT REFERENCES Artists(id) ON DELETE CASCADE NOT NULL,
    track_id  INT REFERENCES Tracks(id)  ON DELETE CASCADE NOT NULL,

    PRIMARY KEY(artist_id, track_id)
);
