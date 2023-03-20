package postgresql

import (
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"

	db "github.com/go-park-mail-ru/2023_1_Technokaif/init/db/postgresql"
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

func (p *PostgreSQL) GetByID(userID uint32) (models.User, error) {
	query := fmt.Sprintf(
		`SELECT id, 
				version, 
				username, 
				email, 
				password, 
				salt, 
				first_name, 
				last_name, 
				sex, 
				birth_date, 
				avatar_src
		FROM %s 
		WHERE id = $1;`,
		db.PostgresTables.Users)

	var user models.User
	err := p.db.Get(&user, query, userID)
	if err == sql.ErrNoRows {
		return models.User{}, fmt.Errorf("(repo) %v: %w", err, &models.NoSuchUserError{})
	} else if err != nil {
		return models.User{}, fmt.Errorf("(repo) failed to exec query: %w", err)
	}

	return user, nil
}
