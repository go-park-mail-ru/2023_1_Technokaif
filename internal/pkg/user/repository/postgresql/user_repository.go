package postgresql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"

	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/models"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/user"
)

// PostgreSQL implements user.Repository
type PostgreSQL struct {
	db     *sqlx.DB
	tables user.Tables
}

func NewPostgreSQL(db *sqlx.DB, t user.Tables) *PostgreSQL {
	return &PostgreSQL{
		db:     db,
		tables: t,
	}
}

const errorUserExists = "unique_violation"

func (p *PostgreSQL) Check(ctx context.Context, userID uint32) error {
	query := fmt.Sprintf(
		`SELECT EXISTS(
			SELECT id
			FROM %s
			WHERE id = $1
		);`,
		p.tables.Users())

	var exists bool
	err := p.db.Get(&exists, query, userID)
	if err != nil {
		return fmt.Errorf("(repo) failed to exec query: %w", err)
	}

	if !exists {
		return fmt.Errorf("(repo) %w: %w", &models.NoSuchUserError{UserID: userID}, err)
	}

	return nil
}

func (p *PostgreSQL) GetByID(ctx context.Context, userID uint32) (*models.User, error) {
	query := fmt.Sprintf(
		`SELECT id, 
				version, 
				username, 
				email, 
				password_hash, 
				salt, 
				first_name, 
				last_name, 
				birth_date, 
				avatar_src
		FROM %s 
		WHERE id = $1;`,
		p.tables.Users())

	row := p.db.QueryRowContext(ctx, query, userID)
	var u models.User
	err := row.Scan(&u.ID, &u.Version, &u.Username, &u.Email, &u.Password, &u.Salt,
		&u.FirstName, &u.LastName, &u.BirthDate.Time, &u.AvatarSrc)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &models.User{}, fmt.Errorf("(repo) %w: %v", &models.NoSuchUserError{UserID: userID}, err)
		}

		return &models.User{}, fmt.Errorf("(repo) failed to exec query: %w", err)
	}

	return &u, nil
}

func (p *PostgreSQL) CreateUser(ctx context.Context, u models.User) (uint32, error) {
	query := fmt.Sprintf(
		`INSERT INTO %s 
			(username, email, password_hash, salt, first_name, last_name, birth_date, avatar_src) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id;`,
		p.tables.Users())

	row := p.db.QueryRowContext(ctx, query, u.Username, u.Email, u.Password, u.Salt,
		u.FirstName, u.LastName, u.BirthDate.Format(time.RFC3339), u.AvatarSrc)

	var id uint32

	err := row.Scan(&id)
	if err != nil {
		if pqerr, ok := err.(*pq.Error); ok {
			if pqerr.Code.Name() == errorUserExists {
				return 0, fmt.Errorf("(repo) %w: %v", &models.UserAlreadyExistsError{}, err)
			}
		}

		return id, fmt.Errorf("(Repo) failed to scan from query: %w", err)
	}

	return id, nil
}

func (p *PostgreSQL) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	query := fmt.Sprintf(
		`SELECT id, version, username, email, password_hash, salt, 
			first_name, last_name, birth_date, avatar_src
		FROM %s WHERE (username=$1 OR email=$1);`,
		p.tables.Users())
	row := p.db.QueryRowContext(ctx, query, username)

	var u models.User
	err := row.Scan(&u.ID, &u.Version, &u.Username, &u.Email, &u.Password, &u.Salt,
		&u.FirstName, &u.LastName, &u.BirthDate.Time, &u.AvatarSrc)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("(repo) %w, %v", &models.NoSuchUserError{}, err)
		}

		return nil, fmt.Errorf("(repo) failed to scan from query: %w", err)
	}

	return &u, nil
}

func (p *PostgreSQL) UpdateInfo(ctx context.Context, u *models.User) error {
	query := fmt.Sprintf(
		`UPDATE %s
		SET email = $2,
			first_name = $3,
			last_name = $4,
			birth_date = $5
		WHERE id = $1;`,
		p.tables.Users())
	if _, err := p.db.ExecContext(ctx, query, u.ID, u.Email, u.FirstName, u.LastName,
		u.BirthDate.Format(time.RFC3339)); err != nil {

		return fmt.Errorf("(repo) failed to exec query: %w", err)
	}

	return nil
}

func (p *PostgreSQL) UpdateAvatarSrc(ctx context.Context, userID uint32, avatarSrc string) error {
	query := fmt.Sprintf(
		`UPDATE %s
		SET avatar_src = $2
		WHERE id = $1;`,
		p.tables.Users())
	if _, err := p.db.ExecContext(ctx, query, userID, avatarSrc); err != nil {

		return fmt.Errorf("(repo) failed to exec query: %w", err)
	}

	return nil
}

func (p *PostgreSQL) GetByPlaylist(ctx context.Context, playlistID uint32) ([]models.User, error) {
	query := fmt.Sprintf(
		`SELECT id,
				username,
				email,
				first_name,
				last_name,
				birth_date,
				avatar_src
		FROM %s u
			INNER JOIN %s up ON u.ID = up.user_id
		WHERE up.playlist_id = $1;`,
		p.tables.Users(), p.tables.UsersPlaylists())

	rows, err := p.db.QueryContext(ctx, query, playlistID)
	if err != nil {
		return nil, fmt.Errorf("(repo) failed to exec query: %w", err)
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var u models.User
		err := rows.Scan(&u.ID, &u.Username, &u.Email, &u.FirstName,
			&u.LastName, &u.BirthDate.Time, &u.AvatarSrc)

		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, fmt.Errorf("(repo) %w: %w", &models.NoSuchPlaylistError{PlaylistID: playlistID}, err)
			}

			return nil, fmt.Errorf("(repo) failed to exec query: %w", err)
		}

		users = append(users, u)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("(repo) failed to exec query: %w", err)
	}

	return users, nil
}
