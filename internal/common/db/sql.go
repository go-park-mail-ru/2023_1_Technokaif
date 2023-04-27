package db

import (
	"database/sql"
	"fmt"
)

func CheckTransaction(tx *sql.Tx, repoError *error) {
	if *repoError != nil {
		if err := tx.Rollback(); err != nil {
			*repoError = fmt.Errorf("(repo) failed to rollback transaction: %w: %w", err, *repoError)
		}
	} else {
		if err := tx.Commit(); err != nil {
			*repoError = fmt.Errorf("(repo) failed to commit transaction: %w: %w", err, *repoError)
		}
	}
}
