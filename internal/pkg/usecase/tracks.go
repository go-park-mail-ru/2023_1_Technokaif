package usecase

import (
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/logger"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/models"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/repository"
)

type TrackUsecase struct {
	repo   repository.Track
	logger logger.Logger
}

func NewTrackUsecase(rt repository.Track, l logger.Logger) *TrackUsecase {
	return &TrackUsecase{repo: rt, logger: l}
}

// GetFeed returns tracks for main page
func (t *TrackUsecase) GetFeed() ([]models.TrackFeed, error) {
	tracks, err := t.repo.GetFeed()

	return tracks, err
}
