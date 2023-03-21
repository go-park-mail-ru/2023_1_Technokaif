package postgresql

import (
	"database/sql"
	"fmt"
	"errors"

	"github.com/jmoiron/sqlx"

	db "github.com/go-park-mail-ru/2023_1_Technokaif/init/db/postgresql"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/models"
	"github.com/go-park-mail-ru/2023_1_Technokaif/pkg/logger"
)

// PostgreSQL implements auth.Repository
type PostgreSQL struct {
	db     *sqlx.DB
	logger logger.Logger
}

func NewPostgreSQL(db *sqlx.DB, l logger.Logger) *PostgreSQL {
	return &PostgreSQL{db: db, logger: l}
}

func (p *PostgreSQL) GetUserByAuthData(userID, userVersion uint32) (*models.User, error) {
	query := fmt.Sprintf(
		`SELECT id, version, username, email, password_hash, salt, 
			first_name, last_name, sex, birth_date 
		FROM %s WHERE id=$1 AND version=$2;`, db.PostgresTables.Users)
	row := p.db.QueryRow(query, userID, userVersion)

	var u models.User
	err := row.Scan(&u.ID, &u.Version, &u.Username, &u.Email, &u.Password, &u.Salt,
		&u.FirstName, &u.LastName, &u.Sex, &u.BirhDate.Time)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("(repo) %w: %v", &models.NoSuchUserError{}, err)
	} else if err != nil {
		return nil, fmt.Errorf("(repo) failed to scan from query: %w", err)
	}

	return &u, nil
}

func (p *PostgreSQL) IncreaseUserVersion(userID uint32) error {
	query := fmt.Sprintf(
		`UPDATE %s SET version = version + 1 WHERE id=$1 RETURNING id;`,
		db.PostgresTables.Users)
	row := p.db.QueryRow(query, userID)

	err := row.Scan(&userID)

	if errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("(repo) %w: %v", &models.NoSuchUserError{}, err)
	} else if err != nil {
		return fmt.Errorf("(repo) failed to scan from query: %w", err)
	}

	return nil
}
