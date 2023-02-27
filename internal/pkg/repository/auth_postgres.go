package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/lib/pq"

	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/models"
)

type AuthPostgres struct {
	db *sql.DB
}

func NewAuthPostgres(db *sql.DB) *AuthPostgres {
	return &AuthPostgres{db: db}
}

const ERROR_NAME_USER_EXISTS = "unique_violation"

func (ap *AuthPostgres) CreateUser(u models.User) (int, error) {
	query := fmt.Sprintf("INSERT INTO %s"+
		"(username, email, password_hash, first_name, last_name, user_sex, birth_date)"+
		"VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id;", UsersTable)

	row := ap.db.QueryRow(query,
		u.Username, u.Email, u.Password, u.FirstName, u.LastName, u.Sex, u.BirhDate.Format(time.RFC3339))

	var id int
	err := row.Scan(&id)

	if pqerr, ok := err.(*pq.Error); ok {
		if pqerr.Code.Name() == ERROR_NAME_USER_EXISTS {
			return 0, &UserAlreadyExistsError{}
		} else {
			return 0, err
		}
	}

	return id, err
}

func (ap *AuthPostgres) GetUser(username, password string) (models.User, error) {
	query := fmt.Sprintf("SELECT id, username, email, password_hash, first_name, last_name, user_sex, birth_date "+
		"FROM %s WHERE (username=$1 OR email=$1) AND password_hash=$2;", UsersTable)
	row := ap.db.QueryRow(query, username, password)

	var u models.User
	err := row.Scan(&u.ID, &u.Username, &u.Email, &u.Password,
		&u.FirstName, &u.LastName, &u.Sex, &u.BirhDate.Time)

	if err == sql.ErrNoRows {
		return u, &NoSuchUserError{}
	}

	return u, err
}
