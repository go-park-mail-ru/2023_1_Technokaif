package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/lib/pq"
	"github.com/pkg/errors"

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

const errorUserExists = "unique_violation"

func (ap *AuthPostgres) CreateUser(u models.User) (uint32, error) {
	query := fmt.Sprintf(`INSERT INTO %s 
	(username, email, password_hash, salt, first_name, last_name, sex, birth_date) 
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id;`, usersTable)

	row := ap.db.QueryRow(query, u.Username, u.Email, u.Password, u.Salt,
		u.FirstName, u.LastName, u.Sex, u.BirhDate.Format(time.RFC3339))

	var id uint32

	err := row.Scan(&id)
	if err != nil {
		if pqerr, ok := err.(*pq.Error); ok {
			if pqerr.Code.Name() == errorUserExists {
				return 0, errors.Wrapf(&models.UserAlreadyExistsError{}, "(Repo) %s", err.Error())
			}
		}
	}

	return id, errors.Wrap(err, "(Repo) failed to scan from query")
}

// TODO make helping func for GetUser*
func (ap *AuthPostgres) GetUserByUsername(username string) (*models.User, error) {
	query := fmt.Sprintf(`SELECT id, version, username, email, password_hash, salt, 
		first_name, last_name, sex, birth_date 
		FROM %s WHERE (username=$1 OR email=$1);`, usersTable)
	row := ap.db.QueryRow(query, username)

	var u models.User
	err := row.Scan(&u.ID, &u.Version, &u.Username, &u.Email, &u.Password, &u.Salt,
		&u.FirstName, &u.LastName, &u.Sex, &u.BirhDate.Time)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, errors.Wrapf(&models.NoSuchUserError{}, "(Repo) %s", err.Error())
	}

	return &u, errors.Wrap(err, "(Repo) failed to scan from query")
}

func (ap *AuthPostgres) GetUserByAuthData(userID, userVersion uint32) (*models.User, error) {
	query := fmt.Sprintf(`SELECT id, version, username, email, password_hash, salt, 
		first_name, last_name, sex, birth_date 
		FROM %s WHERE id=$1 AND version=$2;`, usersTable)
	row := ap.db.QueryRow(query, userID, userVersion)

	var u models.User
	err := row.Scan(&u.ID, &u.Version, &u.Username, &u.Email, &u.Password, &u.Salt,
		&u.FirstName, &u.LastName, &u.Sex, &u.BirhDate.Time)

	if errors.Is(err, sql.ErrNoRows) {
		return &u, errors.Wrapf(&models.NoSuchUserError{}, "(Repo) %s", err.Error())
	}

	return &u, errors.Wrap(err, "(Repo) failed to scan from query")
}

func (ap *AuthPostgres) IncreaseUserVersion(userID uint32) error {
	query := fmt.Sprintf("UPDATE %s SET version = version + 1 WHERE id=$1 RETURNING id;", usersTable)
	row := ap.db.QueryRow(query, userID)

	err := row.Scan(&userID)

	if errors.Is(err, sql.ErrNoRows) {
		return errors.Wrapf(&models.NoSuchUserError{}, "(Repo) %s", err.Error())
	}

	return errors.Wrap(err, "(Repo) failed to scan from query")
}
