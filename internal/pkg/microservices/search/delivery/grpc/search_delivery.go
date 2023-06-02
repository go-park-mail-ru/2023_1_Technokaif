package grpc

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	proto "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/microservices/search/proto/generated"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/search"
	"github.com/go-park-mail-ru/2023_1_Technokaif/pkg/logger"
)

type searchGRPC struct {
	searchServices search.Usecase
	logger         logger.Logger

	proto.UnimplementedSearchServer
}

func NewSearchGRPC(searchServices search.Usecase, l logger.Logger) proto.SearchServer {
	return &searchGRPC{
		searchServices: searchServices,
		logger:         l,
	}
}

func (s *searchGRPC) FindAlbums(msg *proto.SearchMsg, stream proto.Search_FindAlbumsServer) error {
	albums, err := s.searchServices.FindAlbums(stream.Context(), msg.Query, msg.Amount)
	if err != nil {
		return status.Error(codes.Internal, err.Error())
	}

	for _, album := range albums {
		resp := &proto.AlbumResponse{
			Id:          album.ID,
			Name:        album.Name,
			Description: nilCheckString(album.Description),
			CoverSrc:    album.CoverSrc,
		}

		if err := stream.Send(resp); err != nil {
			return status.Error(codes.Internal, err.Error())
		}
	}

	return nil
}

func (s *searchGRPC) FindTracks(msg *proto.SearchMsg, stream proto.Search_FindTracksServer) error {
	tracks, err := s.searchServices.FindTracks(stream.Context(), msg.Query, msg.Amount)
	if err != nil {
		return status.Error(codes.Internal, err.Error())
	}

	for _, track := range tracks {
		resp := &proto.TrackResponse{
			Id:            track.ID,
			Name:          track.Name,
			AlbumID:       nilCheckUint32(track.AlbumID),
			AlbumPosition: nilCheckUint32(track.AlbumPosition),
			CoverSrc:      track.CoverSrc,
			RecordSrc:     track.RecordSrc,
			Duration:      track.Duration,
			Listens:       track.Listens,
		}

		if err := stream.Send(resp); err != nil {
			return status.Error(codes.Internal, err.Error())
		}
	}

	return nil
}

func (s *searchGRPC) FindArtists(msg *proto.SearchMsg, stream proto.Search_FindArtistsServer) error {
	artists, err := s.searchServices.FindArtists(stream.Context(), msg.Query, msg.Amount)
	if err != nil {
		return status.Error(codes.Internal, err.Error())
	}

	for _, artist := range artists {
		resp := &proto.ArtistResponse{
			Id:        artist.ID,
			UserID:    nilCheckUint32(artist.UserID),
			Name:      artist.Name,
			AvatarSrc: artist.AvatarSrc,
			Listens:   artist.Listens,
		}

		if err := stream.Send(resp); err != nil {
			return status.Error(codes.Internal, err.Error())
		}
	}

	return nil
}

func (s *searchGRPC) FindPlaylists(msg *proto.SearchMsg, stream proto.Search_FindPlaylistsServer) error {
	playlists, err := s.searchServices.FindPlaylists(stream.Context(), msg.Query, msg.Amount)
	if err != nil {
		return status.Error(codes.Internal, err.Error())
	}

	for _, playlist := range playlists {
		resp := &proto.PlaylistResponse{
			Id:          playlist.ID,
			Name:        playlist.Name,
			Description: nilCheckString(playlist.Description),
			CoverSrc:    playlist.CoverSrc,
		}

		if err := stream.Send(resp); err != nil {
			return status.Error(codes.Internal, err.Error())
		}
	}

	return nil
}

func nilCheckUint32(val *uint32) uint32 {
	if val == nil {
		return 0
	}
	return *val
}

func nilCheckString(val *string) string {
	if val == nil {
		return ""
	}
	return *val
}
