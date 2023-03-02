package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/lib/pq"

	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/logger"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/models"
)

// AuthPostgres implements Auth
type AuthPostgres struct {
	db     *sql.DB
	logger logger.Logger
}

func NewAuthPostgres(db *sql.DB, l logger.Logger) *AuthPostgres {
	return &AuthPostgres{db: db, logger: l}
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
	if err != nil {
		ap.logger.Error(err.Error())
	} // logger

	if pqerr, ok := err.(*pq.Error); ok {
		if pqerr.Code.Name() == ERROR_NAME_USER_EXISTS {
			return 0, &UserAlreadyExistsError{}
		} else {
			return 0, err
		}
	}

	return id, err
}

// TODO make helping func for GetUser*
func (ap *AuthPostgres) GetUserByCreds(username, password string) (*models.User, error) {
	query := fmt.Sprintf("SELECT id, user_version, username, email, password_hash, first_name, last_name, user_sex, birth_date "+
		"FROM %s WHERE (username=$1 OR email=$1) AND password_hash=$2;", UsersTable)
	row := ap.db.QueryRow(query, username, password)

	var u models.User
	err := row.Scan(&u.ID, &u.Version, &u.Username, &u.Email, &u.Password,
		&u.FirstName, &u.LastName, &u.Sex, &u.BirhDate.Time)

	if err != nil {
		ap.logger.Error(err.Error())
	} // logger
	if err == sql.ErrNoRows {
		return &u, &NoSuchUserError{}
	}

	return &u, err
}

func (ap *AuthPostgres) GetUserByAuthData(userID, userVersion uint) (*models.User, error) {
	query := fmt.Sprintf("SELECT id, user_version, username, email, password_hash, first_name, last_name, user_sex, birth_date "+
		"FROM %s WHERE id=$1 AND user_version=$2;", UsersTable)
	row := ap.db.QueryRow(query, userID, userVersion)

	var u models.User
	err := row.Scan(&u.ID, &u.Version, &u.Username, &u.Email, &u.Password,
		&u.FirstName, &u.LastName, &u.Sex, &u.BirhDate.Time)

	if err != nil {
		ap.logger.Error(err.Error())
	} // logger
	if err == sql.ErrNoRows {
		return &u, &NoSuchUserError{}
	}

	return &u, err
}
