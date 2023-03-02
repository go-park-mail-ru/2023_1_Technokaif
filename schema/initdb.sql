CREATE TYPE SEX AS ENUM('M', 'F', 'O'); 
CREATE TABLE Users
(
    id            SERIAL       PRIMARY KEY,
    version       INT                      NOT NULL DEFAULT 1,
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
    cover_src  TEXT,
    record_src TEXT
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

/* Mock */
INSERT INTO Artists (name) VALUES ('Oxxxymiron');
INSERT INTO Artists (name) VALUES ('Instasamka');

INSERT INTO Albums (name, description) VALUES ('Горгород', 'Антиутопия классная');
INSERT INTO Albums (name, description) VALUES ('Долгий путь домой', 'Качает жёстко');
INSERT INTO Albums (name, description) VALUES ('POPSTAR', 'КАЙФ');

INSERT INTO Tracks (name, album_id) VALUES ('Где нас нет', 1);
INSERT INTO Tracks (name, album_id) VALUES ('Признаки жизни', 2);
INSERT INTO Tracks (name, album_id) VALUES ('За деньги да', 3);
INSERT INTO Tracks (name) VALUES ('Город под подошвой');

INSERT INTO Artists_Albums (artist_id, album_id) VALUES (1, 1);
INSERT INTO Artists_Albums (artist_id, album_id) VALUES (2, 1);
INSERT INTO Artists_Albums (artist_id, album_id) VALUES (2, 3);

INSERT INTO Artists_Tracks (artist_id, track_id) VALUES (1, 1);
INSERT INTO Artists_Tracks (artist_id, track_id) VALUES (1, 2);
INSERT INTO Artists_Tracks (artist_id, track_id) VALUES (2, 3);
