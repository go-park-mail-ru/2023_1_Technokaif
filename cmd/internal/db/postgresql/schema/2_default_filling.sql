INSERT INTO Artists (name, avatar_src)
VALUES ('Oxxxymiron', '/artists/oxxxymiron.jpg'),
        ('SALUKI', '/artists/saluki.jpg'),
        ('Инстасамка', '/artists/instasamka.png'),
        ('ANIKV', '/artists/anikv.png');

INSERT INTO Albums (name, description, cover_src)
VALUES ('Горгород', 'Антиутопия', '/albums/gorgorod.jpg'),
        ('Властелин Калек', 'Стильная душевная музыка', '/albums/vlkol.jpg');

INSERT INTO Tracks (name, album_id, album_position, cover_src, record_src, duration)
VALUES ('Не с начала', 1, 1, '/tracks/gorgorod.jpg', '/records/1.mp3', 125),
        ('Кем ты стал', 1, 2, '/tracks/gorgorod.jpg', '/records/2.mp3', 236),
        ('Всего лишь писатель', 1, 3, '/tracks/gorgorod.jpg', '/records/3.mp3', 209),
        ('Девочка Пи*дец', 1, 4, '/tracks/gorgorod.jpg', '/records/4.mp3', 163),
        ('Переплетено', 1, 5, '/tracks/gorgorod.jpg', '/records/5.mp3', 291),
        ('Колыбельная', 1, 6, '/tracks/gorgorod.jpg', '/records/6.mp3', 197),
        ('Полигон', 1, 7, '/tracks/gorgorod.jpg', '/records/7.mp3', 220),
        ('Накануне', 1, 8, '/tracks/gorgorod.jpg', '/records/8.mp3', 221),
        ('Слово мэра', 1, 9, '/tracks/gorgorod.jpg', '/records/9.mp3', 240),
        ('Башня из слоновой кости', 1, 10, '/tracks/gorgorod.jpg', '/records/10.mp3', 204),
        ('Где нас нет', 1, 11, '/tracks/gorgorod.jpg', '/records/11.mp3', 265),

        ('Дамбо', 2, 1, '/tracks/vlkol.jpg', '/records/12.mp3', 180),
        ('Властелин Калек', 2, 2, '/tracks/vlkol.jpg', '/records/13.mp3', 163),
        ('Поломка', 2, 3, '/tracks/vlkol.jpg', '/records/14.mp3', 130),
        ('Бензобак', 2, 4, '/tracks/vlkol.jpg', '/records/15.mp3', 122),
        ('Пьяный весь район', 2, 5, '/tracks/vlkol.jpg', '/records/16.mp3', 150),
        ('Тупик', 2, 6, '/tracks/vlkol.jpg', '/records/17.mp3', 184),
        ('Пекло', 2, 7, '/tracks/vlkol.jpg', '/records/18.mp3', 148),
        ('ВЛАДИВОСТОК 3000', 2, 8, '/tracks/vlkol.jpg', '/records/19.mp3', 206),
        ('Болевой шок', 2, 9, '/tracks/vlkol.jpg', '/records/20.mp3', 218),
        ('Алый', 2, 10, '/tracks/vlkol.jpg', '/records/21.mp3', 138),
        ('Понт', 2, 11, '/tracks/vlkol.jpg', '/records/22.mp3', 236),
        ('Решето', 2, 12, '/tracks/vlkol.jpg', '/records/23.mp3', 193),
        ('NNN705', 2, 13, '/tracks/vlkol.jpg', '/records/24.mp3', 145),
        ('Улыбка', 2, 14, '/tracks/vlkol.jpg', '/records/25.mp3', 254),
        ('Ilford XP2 Super', 2, 15, '/tracks/vlkol.jpg', '/records/26.mp3', 204);

INSERT INTO Tracks (name, cover_src, record_src, duration)
VALUES ('LAGG OUT', '/tracks/laggOut.jpeg', '/records/27.mp3', 279),
       ('Город под подошвой', '/tracks/gorodPodPod.png', '/records/28.mp3', 229),
       ('За деньги да', '/tracks/zadengida.png', '/records/29.mp3', 119),
       ('Mommy', '/tracks/mommy.png', '/records/30.mp3', 93),
       ('Juicy', '/tracks/juicy.png', '/records/31.mp3', 123);

INSERT INTO Artists_Albums (artist_id, album_id)
VALUES (1, 1),
       (2, 2);

INSERT INTO Artists_Tracks (artist_id, track_id)
VALUES (1, 1), (1, 2), (1, 3), (1, 4), (1, 5), (1, 6), (1, 7), (1, 8), (1, 9), (1, 10),
        (1, 11), (2, 12), (2, 13), (2, 14), (2, 15), (2, 16), (2, 17), (2, 18), (2, 19), (2, 20),
        (2, 21), (4, 21), (2, 22), (4, 22), (2, 23), (2, 24), (2, 25), (2, 26), (2, 27), (1, 28), (3, 29), (3, 30),
        (3, 31);
