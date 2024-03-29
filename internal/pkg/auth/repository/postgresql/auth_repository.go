package postgresql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/models"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/auth"
)

// PostgreSQL implements auth.Repository
type PostgreSQL struct {
	db     *sqlx.DB
	tables auth.Tables
}

func NewPostgreSQL(db *sqlx.DB, t auth.Tables) *PostgreSQL {
	return &PostgreSQL{
		db:     db,
		tables: t,
	}
}

func (p *PostgreSQL) GetUserByAuthData(ctx context.Context, userID, userVersion uint32) (*models.User, error) {
	query := fmt.Sprintf(
		`SELECT id, version, username, email, password_hash, salt, 
			first_name, last_name, birth_date, avatar_src 
		FROM %s
		WHERE id = $1 AND version = $2;`,
		p.tables.Users())
	row := p.db.QueryRowContext(ctx, query, userID, userVersion)

	var u models.User
	err := row.Scan(&u.ID, &u.Version, &u.Username, &u.Email, &u.Password, &u.Salt,
		&u.FirstName, &u.LastName, &u.BirthDate.Time, &u.AvatarSrc)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("(repo) %w: %w", &models.NoSuchUserError{UserID: userID}, err)
		}

		return nil, fmt.Errorf("(repo) failed to scan from query: %w", err)
	}

	return &u, nil
}

func (p *PostgreSQL) IncreaseUserVersion(ctx context.Context, userID uint32) error {
	query := fmt.Sprintf(
		`UPDATE %s
		SET version = version + 1
		WHERE id = $1
		RETURNING id;`,
		p.tables.Users())
	row := p.db.QueryRowContext(ctx, query, userID)

	err := row.Scan(&userID)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("(repo) %w: %w", &models.NoSuchUserError{UserID: userID}, err)
		}

		return fmt.Errorf("(repo) failed to scan from query: %w", err)
	}

	return nil
}

func (p *PostgreSQL) UpdatePassword(ctx context.Context, userID uint32, passwordHash, salt string) error {
	query := fmt.Sprintf(
		`UPDATE %s
		SET password_hash = $1,
			salt = $2
		WHERE id = $3
		RETURNING id;`,
		p.tables.Users())
	row := p.db.QueryRowContext(ctx, query, passwordHash, salt, userID)

	if err := row.Scan(&userID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("(repo) %w: %v", &models.NoSuchUserError{}, err)
		}

		return fmt.Errorf("(repo) failed to scan from query: %w", err)
	}

	return nil
}
