package postgresql

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"

	dbInit "github.com/go-park-mail-ru/2023_1_Technokaif/init/db/postgresql"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/models"
	"github.com/go-park-mail-ru/2023_1_Technokaif/pkg/logger"
)

// PostgreSQL implements user.Repository
type PostgreSQL struct {
	db     *sqlx.DB
	logger logger.Logger
}

func NewPostgreSQL(db *sqlx.DB, l logger.Logger) *PostgreSQL {
	return &PostgreSQL{db: db, logger: l}
}

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
		dbInit.PostgresTables.Users)

	
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
