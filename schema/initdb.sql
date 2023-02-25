CREATE TYPE SEX AS ENUM('M', 'F', 'O'); 
CREATE TABLE Users
(
    id SERIAL PRIMARY KEY,
    username      VARCHAR(20)  UNIQUE NOT NULL,
    email         VARCHAR(30)  UNIQUE NOT NULL,
    password_hash VARCHAR(256)        NOT NULL,
    first_name    VARCHAR(20)         NOT NULL,
    last_name     VARCHAR(20)         NOT NULL,
    user_sex      SEX                 NOT NULL
);
