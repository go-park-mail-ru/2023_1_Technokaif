package postgresql

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"

	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/models"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/user"
	"github.com/go-park-mail-ru/2023_1_Technokaif/pkg/logger"
)

// PostgreSQL implements user.Repository
type PostgreSQL struct {
	db     *sqlx.DB
	tables user.Tables
	logger logger.Logger
}

func NewPostgreSQL(db *sqlx.DB, t user.Tables, l logger.Logger) *PostgreSQL {
	return &PostgreSQL{
		db:     db,
		tables: t,
		logger: l,
	}
}

const errorUserExists = "unique_violation"

func (p *PostgreSQL) GetByID(userID uint32) (*models.User, error) {
	query := fmt.Sprintf(
		`SELECT id, 
				version, 
				username, 
				email, 
				password_hash, 
				salt, 
				first_name, 
				last_name, 
				sex, 
				birth_date, 
				COALESCE(avatar_src, '')
		FROM %s 
		WHERE id = $1;`,
		p.tables.Users())

	row := p.db.QueryRow(query, userID)
	var u models.User
	err := row.Scan(&u.ID, &u.Version, &u.Username, &u.Email, &u.Password, &u.Salt,
		&u.FirstName, &u.LastName, &u.Sex, &u.BirhDate.Time, &u.AvatarSrc)

	if errors.Is(err, sql.ErrNoRows) {
		return &models.User{}, fmt.Errorf("(repo) %w: %v", &models.NoSuchUserError{}, err)
	} else if err != nil {
		return &models.User{}, fmt.Errorf("(repo) failed to exec query: %w", err)
	}

	return &u, nil
}

func (p *PostgreSQL) CreateUser(u models.User) (uint32, error) {
	query := fmt.Sprintf(
		`INSERT INTO %s 
			(username, email, password_hash, salt, first_name, last_name, sex, birth_date) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id;`,
		p.tables.Users())

	row := p.db.QueryRow(query, u.Username, u.Email, u.Password, u.Salt,
		u.FirstName, u.LastName, u.Sex, u.BirhDate.Format(time.RFC3339))

	var id uint32

	err := row.Scan(&id)
	if err != nil {
		if pqerr, ok := err.(*pq.Error); ok {
			if pqerr.Code.Name() == errorUserExists {
				return 0, fmt.Errorf("(repo) %w: %v", &models.UserAlreadyExistsError{}, err)
			}
		} else {
			return id, fmt.Errorf("(Repo) failed to scan from query: %w", err)
		}
	}

	return id, nil
}

// TODO make helping func for GetUser*
func (p *PostgreSQL) GetUserByUsername(username string) (*models.User, error) {
	query := fmt.Sprintf(
		`SELECT id, version, username, email, password_hash, salt, 
			first_name, last_name, sex, birth_date 
		FROM %s WHERE (username=$1 OR email=$1);`,
		p.tables.Users())
	row := p.db.QueryRow(query, username)

	var u models.User
	err := row.Scan(&u.ID, &u.Version, &u.Username, &u.Email, &u.Password, &u.Salt,
		&u.FirstName, &u.LastName, &u.Sex, &u.BirhDate.Time)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("(repo) %w, %v", &models.NoSuchUserError{}, err)
	} else if err != nil {
		return nil, fmt.Errorf("(repo) failed to scan from query: %w", err)
	}

	return &u, nil
}
