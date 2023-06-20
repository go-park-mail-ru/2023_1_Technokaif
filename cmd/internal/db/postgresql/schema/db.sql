-- Таблица пользователей приложения
-- Удовлетворяет НФ Бойса-Кодда
CREATE TABLE Users
(
    id            SERIAL       PRIMARY KEY,
    version       INT                      NOT NULL DEFAULT 1,
    username      VARCHAR(20)  UNIQUE      NOT NULL, -- У каждого пользователя свой уникальный username
    email         VARCHAR(255) UNIQUE      NOT NULL, -- у каждого пользователя свой уникальный почтовый адрес
    password_hash VARCHAR(256)             NOT NULL,
    salt          VARCHAR(64)              NOT NULL,
    first_name    VARCHAR(20)              NOT NULL,
    last_name     VARCHAR(20)              NOT NULL,
    birth_date    DATE                     NOT NULL,
    avatar_src    TEXT
);

-- Таблица языков для Full Text Search по сущностям
CREATE TABLE Dict_langs
(
    id   SERIAL    PRIMARY KEY,
    lang REGCONFIG             NOT NULL
);

-- Таблица артистов (музыканты, группы)
-- Удовлетворяет 2 НФ, т.к. lang_id определяется значением name
CREATE TABLE Artists
(
    id         SERIAL      PRIMARY KEY,
    name       VARCHAR(30)                                              NOT NULL,
    avatar_src TEXT                                                     NOT NULL,
    lang_id    INT         REFERENCES Dict_langs(id) ON DELETE SET NULL
);

-- Таблица музыкальных альбомов
-- Удовлетворяет 2 НФ, т.к. lang_id определяется значением name
CREATE TABLE Albums
(
    id          SERIAL        PRIMARY KEY,
    name        VARCHAR(40)               NOT NULL,
    description VARCHAR(2000),
    cover_src   TEXT                      NOT NULL,
    lang_id    INT         REFERENCES Dict_langs(id) ON DELETE SET NULL
);

-- Таблица связи М:М артистов и альбомов
CREATE TABLE Artists_Albums
(
    artist_id INT REFERENCES Artists(id) ON DELETE CASCADE NOT NULL,
    album_id  INT REFERENCES Albums(id)  ON DELETE CASCADE NOT NULL,

    PRIMARY KEY(artist_id, album_id) -- Составной PK (unique + not null)
);

-- Таблица треков (музыкальных композиций)
-- Удовлетворяет 2 НФ, т.к. lang_id определяется значением name
CREATE TABLE Tracks
(
    id             SERIAL      PRIMARY KEY,
    name           VARCHAR(40)                                          NOT NULL,
    album_id       INT         REFERENCES Albums(id)  ON DELETE CASCADE,
    album_position INT,
    cover_src      TEXT                                                 NOT NULL,
    record_src     TEXT                                                 NOT NULL,
    duration       INT                                                  NOT NULL,
    lang_id        INT         REFERENCES Dict_langs(id) ON DELETE SET NULL,

    -- unique, т.к. внутри одного альбома
    -- два трека не могут находиться на одинаковой позиции
    UNIQUE(album_id, album_position)
);

-- Таблица связи М:М артистов и треков
CREATE TABLE Artists_Tracks
(
    artist_id INT REFERENCES Artists(id) ON DELETE CASCADE NOT NULL,
    track_id  INT REFERENCES Tracks(id)  ON DELETE CASCADE NOT NULL,

    PRIMARY KEY(artist_id, track_id) -- Составной PK (unique + not null)
);

-- Таблица пользовательских плейлистов
-- Удовлетворяет 2 НФ, т.к. lang_id определяется значением name
CREATE TABLE Playlists
(
    id          SERIAL        PRIMARY KEY,
    name        VARCHAR(60)               NOT NULL,
    description VARCHAR(2000),
    cover_src   TEXT,
    lang_id     INT           REFERENCES Dict_langs(id) ON DELETE SET NULL
);

-- Таблица связи М:М пользователей и их плейлистов
CREATE TABLE Users_Playlists
(
    user_id     INT REFERENCES Users(id)     ON DELETE SET NULL,
    playlist_id INT REFERENCES Playlists(id) ON DELETE CASCADE,
    created_at  TIMESTAMPTZ DEFAULT NOW()                       NOT NULL,

    PRIMARY KEY(user_id, playlist_id) -- Составной PK (unique + not null)
);

-- Таблица связи М:М треков и пользовательских плейлистов
CREATE TABLE Playlists_Tracks
(
    playlist_id INT REFERENCES Playlists(id) ON DELETE CASCADE,
    track_id    INT REFERENCES Tracks        ON DELETE CASCADE,
    added_at    TIMESTAMPTZ DEFAULT NOW()                      NOT NULL,

    PRIMARY KEY(playlist_id, track_id) -- Составной PK (unique + not null)
);

-- Таблица связи М:М пользователей и понравившихся им альбомов
CREATE TABLE Liked_albums
(
    user_id   INT REFERENCES Users(id)  ON DELETE CASCADE NOT NULL,
    album_id  INT REFERENCES Albums(id) ON DELETE CASCADE NOT NULL,
    liked_at  TIMESTAMPTZ DEFAULT NOW()                   NOT NULL,

    PRIMARY KEY(user_id, album_id) -- Составной PK (unique + not null)
);

-- Таблица связи М:М пользователей и понравившихся им артистов
CREATE TABLE Liked_artists
(
    user_id    INT REFERENCES Users(id)   ON DELETE CASCADE NOT NULL,
    artist_id  INT REFERENCES Artists(id) ON DELETE CASCADE NOT NULL,
    liked_at   TIMESTAMPTZ DEFAULT NOW()                    NOT NULL,

    PRIMARY KEY(user_id, artist_id) -- Составной PK (unique + not null)
);

-- Таблица связи М:М пользователей и понравившихся им треков
CREATE TABLE Liked_tracks
(
    user_id   INT REFERENCES Users(id)  ON DELETE CASCADE NOT NULL,
    track_id  INT REFERENCES Tracks(id) ON DELETE CASCADE NOT NULL,
    liked_at  TIMESTAMPTZ DEFAULT NOW()                   NOT NULL,

    PRIMARY KEY(user_id, track_id) -- Составной PK (unique + not null)
);

-- Таблица связи М:М пользователей и понравившихся им плейлистов
CREATE TABLE Liked_playlists
(
    user_id     INT REFERENCES Users(id)     ON DELETE CASCADE NOT NULL,
    playlist_id INT REFERENCES Playlists(id) ON DELETE CASCADE NOT NULL,
    liked_at    TIMESTAMPTZ DEFAULT NOW()                      NOT NULL,

    PRIMARY KEY(user_id, playlist_id) -- Составной PK (unique + not null)
);

-- B-tree-индекс для ускорения выборки плейлистов по юзеру и их сортировки
CREATE INDEX idx_btree_playlists_user ON Users_Playlists USING btree (user_id, created_at DESC, playlist_id);

-- B-tree-индекс для ускорения поиска плейлистов по name по совпадению подслов
CREATE INDEX idx_btree_playlists ON Playlists USING btree (LOWER(name) varchar_pattern_ops);
