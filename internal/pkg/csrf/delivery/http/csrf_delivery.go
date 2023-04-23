package delivery

import (
	"net/http"

	commonHttp "github.com/go-park-mail-ru/2023_1_Technokaif/internal/common/http"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/token"
	"github.com/go-park-mail-ru/2023_1_Technokaif/pkg/logger"
)

type Handler struct {
	tokenServices token.Usecase
	logger        logger.Logger
}

func NewHandler(tu token.Usecase, l logger.Logger) *Handler {
	return &Handler{
		tokenServices: tu,
		logger:        l,
	}
}

// swaggermock
func (h *Handler) GetCSRF(w http.ResponseWriter, r *http.Request) {
	user, err := commonHttp.GetUserFromRequest(r)
	if err != nil {
		h.logger.Infof("failed to get CSRF: %s", err.Error())
		commonHttp.ErrorResponse(w, invalidAccessToken, http.StatusUnauthorized, h.logger)
		return
	}

	token, err := h.tokenServices.GenerateCSRFToken(user.ID)
	if err != nil {
		h.logger.Errorf("failed to get CSRF: %s", err.Error())
		commonHttp.ErrorResponse(w, csrfGetError, http.StatusInternalServerError, h.logger)
		return
	}

	resp := getCSRFResponce{CSRF: token}
	commonHttp.SuccessResponse(w, resp, h.logger)
}
