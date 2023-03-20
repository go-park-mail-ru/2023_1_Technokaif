package artist_usecase

import (
	"fmt"

	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/models"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/track"
	"github.com/go-park-mail-ru/2023_1_Technokaif/pkg/logger"
)

type trackUsecase struct {
	repo   track.TrackRepository
	logger logger.Logger
}

func NewTrackUsecase(tr track.TrackRepository, l logger.Logger) track.TrackUsecase {
	return &trackUsecase{repo: tr, logger: l}
}

func (tu *trackUsecase) Create(track models.Track) error {
	if err := tu.repo.Insert(track); err != nil {
		return fmt.Errorf("(usecase) can't create track in repository: %w", err)
	}

	return nil
}

func (tu *trackUsecase) GetByID(trackID uint32) (models.Track, error) {
	track, err := tu.repo.GetByID(trackID)
	if err != nil {
		return models.Track{}, fmt.Errorf("(usecase) can't get track from repository: %w", err)
	}

	return track, nil
}

func (tu *trackUsecase) Change(track models.Track) error {
	if err := tu.repo.Update(track); err != nil {
		return fmt.Errorf("(usecase) can't get track from repository: %w", err)
	}

	return nil
}

func (tu *trackUsecase) Delete(trackID uint32) error {
	if err := tu.repo.Delete(trackID); err != nil {
		return fmt.Errorf("(usecase) can't get track from repository: %w", err)
	}

	return nil
}

func (tu *trackUsecase) GetFeed() ([]models.Track, error) {
	tracks, err := tu.repo.GetFeed()
	if err != nil {
		return nil, fmt.Errorf("(usecase) can't get tracks from repository: %w", err)
	}

	return tracks, nil
}

func (tu *trackUsecase) GetByAlbum(albumID uint32) ([]models.Track, error) {
	tracks, err := tu.repo.GetByAlbum(albumID)
	if err != nil {
		return nil, fmt.Errorf("(usecase) can't get tracks from repository: %w", err)
	}

	return tracks, nil
}

func (tu *trackUsecase) GetByArtist(artistID uint32) ([]models.Track, error) {
	tracks, err := tu.repo.GetByArtist(artistID)
	if err != nil {
		return nil, fmt.Errorf("(usecase) can't get tracks from repository: %w", err)
	}

	return tracks, nil
}

func (tu *trackUsecase) GetLikedByUser(userID uint32) ([]models.Track, error) {
	tracks, err := tu.repo.GetLikedByUser(userID)
	if err != nil {
		return nil, fmt.Errorf("(usecase) can't get tracks from repository: %w", err)
	}

	return tracks, nil
}
