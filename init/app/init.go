package app

import (
	"database/sql"

	"github.com/go-chi/chi"
	"github.com/go-park-mail-ru/2023_1_Technokaif/init/router"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/logger"
)

func Init(db *sql.DB, logger logger.Logger) (*chi.Mux, error) {

	return router.InitRouter()
}
