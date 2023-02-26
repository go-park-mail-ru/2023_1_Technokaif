package usecase

import (
	"github.com/go-park-mail-ru/2023_1_Technokaif/pkg/repository"
)

type TrackUsecase struct {
	repo repository.Track
}

func (t *TrackUsecase) GetTracks() []Track {
	return nil
}

func NewTrackUsecase(ra repository.Track) *TrackUsecase {
	return &TrackUsecase{repo: ra}
}
