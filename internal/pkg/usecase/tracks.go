package usecase

import (
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/models"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/repository"
)

type TrackUsecase struct {
	repo repository.Track
}

func NewTrackUsecase(rt repository.Track) *TrackUsecase {
	return &TrackUsecase{repo: rt}
}

// GetFeed returns tracks for main page
func (t *TrackUsecase) GetFeed() ([]models.TrackFeed, error) {
	tracks, err := t.repo.GetFeed()
	if err != nil {
		return nil, err
	}

	return tracks, nil
}
