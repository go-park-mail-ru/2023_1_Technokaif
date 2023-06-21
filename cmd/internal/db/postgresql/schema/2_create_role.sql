CREATE USER technokaif;

GRANT SELECT, INSERT, UPDATE            ON Users TO technokaif;
GRANT SELECT, INSERT, DELETE            ON Artists TO technokaif;
GRANT SELECT, INSERT, DELETE            ON Albums TO technokaif;
GRANT SELECT, INSERT, DELETE            ON Artists_Albums TO technokaif;
GRANT SELECT, INSERT, DELETE            ON Tracks TO technokaif;
GRANT SELECT, INSERT, DELETE            ON Artists_Tracks TO technokaif;
GRANT SELECT, INSERT, UPDATE, DELETE    ON Playlists TO technokaif;

GRANT SELECT, INSERT, UPDATE, DELETE    ON Users_Playlists TO technokaif;
GRANT SELECT, INSERT, DELETE            ON Playlists_Tracks TO technokaif;
GRANT SELECT, INSERT, DELETE            ON Liked_albums TO technokaif;
GRANT SELECT, INSERT, DELETE            ON Liked_artists TO technokaif;
GRANT SELECT, INSERT, DELETE            ON Liked_albums TO technokaif;
GRANT SELECT, INSERT, DELETE            ON Liked_playlists TO technokaif;
GRANT SELECT                            ON Dict_langs TO technokaif;
