CREATE TYPE SEX AS ENUM('M', 'F', 'O'); 
CREATE TABLE Users
(
    id            SERIAL       PRIMARY KEY,
    version       INT                      NOT NULL DEFAULT 1,
    username      VARCHAR(20)  UNIQUE      NOT NULL,
    email         VARCHAR(30)  UNIQUE      NOT NULL,
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
    name       VARCHAR(30)             NOT NULL,
    avatar_src TEXT                    NOT NULL
);

CREATE TABLE Albums
(
    id          SERIAL        PRIMARY KEY,
    name        VARCHAR(40)               NOT NULL,
    description VARCHAR(2000),
    cover_src   TEXT
);

CREATE TABLE Tracks
(
    id         SERIAL      PRIMARY KEY,
    name       VARCHAR(40)                                          NOT NULL,
    album_id   INT         REFERENCES Albums(id)  ON DELETE CASCADE,
    cover_src  TEXT                                                 NOT NULL,
    record_src TEXT                                                 NOT NULL
);

CREATE TABLE Listens
(
    id       SERIAL PRIMARY KEY,
    user_id  INT    REFERENCES Users(id)  ON DELETE SET NULL,
    track_id INT    REFERENCES Tracks(id) ON DELETE CASCADE  NOT NULL
);

CREATE TABLE Artists_Albums
(
    artist_id INT REFERENCES Artists(id) ON DELETE CASCADE NOT NULL,
    album_id  INT REFERENCES Albums(id)  ON DELETE CASCADE NOT NULL,

    PRIMARY KEY(artist_id, album_id)
);

CREATE TABLE Artists_Tracks
(
    artist_id INT REFERENCES Artists(id) ON DELETE CASCADE NOT NULL,
    track_id  INT REFERENCES Tracks(id)  ON DELETE CASCADE NOT NULL,

    PRIMARY KEY(artist_id, track_id)
);

CREATE TABLE Liked_albums
(
    user_id INT REFERENCES Users(id) ON DELETE CASCADE NOT NULL,
    album_id  INT REFERENCES Albums(id)  ON DELETE CASCADE NOT NULL,

    PRIMARY KEY(user_id, album_id)
);

CREATE TABLE Liked_artists
(
    user_id INT REFERENCES Users(id) ON DELETE CASCADE NOT NULL,
    artist_id  INT REFERENCES Artists(id)  ON DELETE CASCADE NOT NULL,

    PRIMARY KEY(user_id, artist_id)
);

CREATE TABLE Liked_tracks
(
    user_id INT REFERENCES Users(id) ON DELETE CASCADE NOT NULL,
    track_id  INT REFERENCES Tracks(id)  ON DELETE CASCADE NOT NULL,

    PRIMARY KEY(user_id, track_id)
);

/* Default filling */
INSERT INTO Artists (name, avatar_src)
VALUES ('Oxxxymiron', '/artists/oxxxymiron.jpg'),
        ('SALUKI', '/artists/saluki.jpg');

INSERT INTO Albums (name, description, cover_src)
VALUES ('Горгород', 'Антиутопия', '/albums/gorgorod.jpg'),
    ('Долгий путь домой', 'Грайм из Лондона', '/albums/longWayHome.png'),
    ('На Человека', 'Стильная музыка от русского Канье Уэста', '/albums/onHuman.jpg');

INSERT INTO Tracks (name, album_id, cover_src, record_src)
VALUES ('Где нас нет', 1, '/tracks/gorgorod.jpg', '/tracks/gorgorod.wav'),
       ('Признаки жизни', 2, '/tracks/longWayHome.png', '/tracks/longWayHome.mp3');
INSERT INTO Tracks (name, cover_src, record_src)
VALUES ('LAGG OUT', '/tracks/laggOut.jpeg', '/tracks/laggOut.wav'),
       ('Город под подошвой', '/tracks/gorodPodPod.png', '/tracks/gorodPodPod.png');

INSERT INTO Artists_Albums (artist_id, album_id)
VALUES (1, 1),
       (1, 2),
       (2, 3);

INSERT INTO Artists_Tracks (artist_id, track_id)
VALUES (1, 1),
       (1, 2),
       (2, 3),
       (1, 4);
