package grpc

import (
	"context"
	"errors"
	"io"

	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/models"
	proto "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/microservices/search/proto/generated"
)

type SearchAgent struct {
	client proto.SearchClient
}

func NewAuthAgent(c proto.SearchClient) *SearchAgent {
	return &SearchAgent{
		client: c,
	}
}

func (s *SearchAgent) FindAlbums(ctx context.Context, query string, amount uint32) ([]models.Album, error) {
	msg := &proto.SearchMsg{
		Query:  query,
		Amount: amount,
	}

	grpcCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	stream, err := s.client.FindAlbums(grpcCtx, msg)
	if err != nil {
		return nil, err
	}

	albums := make([]models.Album, 0, amount)
	for i := 0; uint32(i) < amount; i++ {
		albumProto, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return nil, err
		}

		album := models.Album{
			ID:          albumProto.Id,
			Name:        albumProto.Name,
			CoverSrc:    albumProto.CoverSrc,
			Description: nilConvertString(albumProto.Description),
		}

		albums = append(albums, album)
	}

	return albums, nil
}

func (s *SearchAgent) FindArtists(ctx context.Context, query string, amount uint32) ([]models.Artist, error) {
	msg := &proto.SearchMsg{
		Query:  query,
		Amount: amount,
	}

	grpcCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	stream, err := s.client.FindArtists(grpcCtx, msg)
	if err != nil {
		return nil, err
	}

	artists := make([]models.Artist, 0, amount)
	for i := 0; uint32(i) < amount; i++ {
		artistProto, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return nil, err
		}

		artist := models.Artist{
			ID:        artistProto.Id,
			UserID:    nilConvertUint32(artistProto.UserID),
			Name:      artistProto.Name,
			AvatarSrc: artistProto.AvatarSrc,
		}

		artists = append(artists, artist)
	}

	return artists, nil
}

func (s *SearchAgent) FindTracks(ctx context.Context, query string, amount uint32) ([]models.Track, error) {
	msg := &proto.SearchMsg{
		Query:  query,
		Amount: amount,
	}

	grpcCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	stream, err := s.client.FindTracks(grpcCtx, msg)
	if err != nil {
		return nil, err
	}

	tracks := make([]models.Track, 0, amount)
	for i := 0; uint32(i) < amount; i++ {
		trackProto, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return nil, err
		}

		track := models.Track{
			ID:            trackProto.Id,
			Name:          trackProto.Name,
			AlbumID:       nilConvertUint32(trackProto.AlbumID),
			AlbumPosition: nilConvertUint32(trackProto.AlbumPosition),
			CoverSrc:      trackProto.CoverSrc,
			RecordSrc:     trackProto.CoverSrc,
			Duration:      trackProto.Duration,
			Listens:       trackProto.Listens,
		}

		tracks = append(tracks, track)
	}

	return tracks, nil
}

func (s *SearchAgent) FindPlaylists(ctx context.Context, query string, amount uint32) ([]models.Playlist, error) {
	msg := &proto.SearchMsg{
		Query:  query,
		Amount: amount,
	}

	grpcCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	stream, err := s.client.FindPlaylists(grpcCtx, msg)
	if err != nil {
		return nil, err
	}

	playlists := make([]models.Playlist, 0, amount)
	for i := 0; uint32(i) < amount; i++ {
		playlistProto, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return nil, err
		}

		playlist := models.Playlist{
			ID:          playlistProto.Id,
			Name:        playlistProto.Name,
			Description: nilConvertString(playlistProto.Description),
			CoverSrc:    playlistProto.CoverSrc,
		}

		playlists = append(playlists, playlist)
	}

	return playlists, nil
}

func nilConvertString(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func nilConvertUint32(val uint32) *uint32 {
	if val == 0 {
		return nil
	}
	return &val
}
