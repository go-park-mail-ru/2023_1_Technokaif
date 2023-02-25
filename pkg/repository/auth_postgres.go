package repository

import (
	"database/sql"
	"fmt"

	"github.com/go-park-mail-ru/2023_1_Technokaif/models"
)

type AuthPostgres struct {
	db *sql.DB
}

func NewAuthPostgres(db *sql.DB) *AuthPostgres {
	return &AuthPostgres{db: db}
}

func (ap *AuthPostgres) CreateUser(u models.User) (int, error) {
	query := fmt.Sprintf("INSERT INTO %s"+
		"(username, email, password_hash, first_name, last_name, user_sex)"+
		"VALUES ($1, $2, $3, $4, $5, $6) RETURNING id;", UsersTable)

	row := ap.db.QueryRow(query,
		u.Username, u.Email, u.Password, u.FirstName, u.LastName, u.Sex)

	var id int
	err := row.Scan(&id)

	return id, err
}

func (ap *AuthPostgres) GetUser(username, password string) (models.User, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE username=$1 AND password_hash=$2;", UsersTable)
	row := ap.db.QueryRow(query, username, password)

	var u models.User
	err := row.Scan(&u.ID, &u.Username, &u.Email, &u.Password,
		&u.FirstName, &u.LastName, &u.Sex)

	return u, err
}
