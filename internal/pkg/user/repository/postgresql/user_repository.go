package user_repository

import (
	"fmt"
	"database/sql"

	"github.com/jmoiron/sqlx"

	"github.com/go-park-mail-ru/2023_1_Technokaif/init/db"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/logger"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/models"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/user"
)

// AuthPostgres implements Auth
type userPostgres struct {
	db     *sqlx.DB
	logger logger.Logger
}

func NewUserPostgres(db *sqlx.DB, l logger.Logger) user.UserRepository {
	return &userPostgres{db: db, logger: l}
}

func (u *userPostgres) GetByID(userID uint32) (models.User, error) {
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
	err := u.db.Get(&user, query, userID); 
	if err == sql.ErrNoRows {
		return models.User{}, fmt.Errorf("(repo) %v: %w", err, &models.NoSuchUserError{})
	} else if err != nil {
		return models.User{}, fmt.Errorf("(repo) failed to exec query: %w", err)
	}

	return user, nil
}