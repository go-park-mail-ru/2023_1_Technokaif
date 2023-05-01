// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.12.4
// source: search.proto

package proto

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// SearchClient is the client API for Search service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type SearchClient interface {
	FindAlbums(ctx context.Context, in *SearchMsg, opts ...grpc.CallOption) (Search_FindAlbumsClient, error)
	FindTracks(ctx context.Context, in *SearchMsg, opts ...grpc.CallOption) (Search_FindTracksClient, error)
	FindPlaylists(ctx context.Context, in *SearchMsg, opts ...grpc.CallOption) (Search_FindPlaylistsClient, error)
	FindArtists(ctx context.Context, in *SearchMsg, opts ...grpc.CallOption) (Search_FindArtistsClient, error)
}

type searchClient struct {
	cc grpc.ClientConnInterface
}

func NewSearchClient(cc grpc.ClientConnInterface) SearchClient {
	return &searchClient{cc}
}

func (c *searchClient) FindAlbums(ctx context.Context, in *SearchMsg, opts ...grpc.CallOption) (Search_FindAlbumsClient, error) {
	stream, err := c.cc.NewStream(ctx, &Search_ServiceDesc.Streams[0], "/auth.Search/FindAlbums", opts...)
	if err != nil {
		return nil, err
	}
	x := &searchFindAlbumsClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Search_FindAlbumsClient interface {
	Recv() (*AlbumResponse, error)
	grpc.ClientStream
}

type searchFindAlbumsClient struct {
	grpc.ClientStream
}

func (x *searchFindAlbumsClient) Recv() (*AlbumResponse, error) {
	m := new(AlbumResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *searchClient) FindTracks(ctx context.Context, in *SearchMsg, opts ...grpc.CallOption) (Search_FindTracksClient, error) {
	stream, err := c.cc.NewStream(ctx, &Search_ServiceDesc.Streams[1], "/auth.Search/FindTracks", opts...)
	if err != nil {
		return nil, err
	}
	x := &searchFindTracksClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Search_FindTracksClient interface {
	Recv() (*TrackResponse, error)
	grpc.ClientStream
}

type searchFindTracksClient struct {
	grpc.ClientStream
}

func (x *searchFindTracksClient) Recv() (*TrackResponse, error) {
	m := new(TrackResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *searchClient) FindPlaylists(ctx context.Context, in *SearchMsg, opts ...grpc.CallOption) (Search_FindPlaylistsClient, error) {
	stream, err := c.cc.NewStream(ctx, &Search_ServiceDesc.Streams[2], "/auth.Search/FindPlaylists", opts...)
	if err != nil {
		return nil, err
	}
	x := &searchFindPlaylistsClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Search_FindPlaylistsClient interface {
	Recv() (*PlaylistResponse, error)
	grpc.ClientStream
}

type searchFindPlaylistsClient struct {
	grpc.ClientStream
}

func (x *searchFindPlaylistsClient) Recv() (*PlaylistResponse, error) {
	m := new(PlaylistResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *searchClient) FindArtists(ctx context.Context, in *SearchMsg, opts ...grpc.CallOption) (Search_FindArtistsClient, error) {
	stream, err := c.cc.NewStream(ctx, &Search_ServiceDesc.Streams[3], "/auth.Search/FindArtists", opts...)
	if err != nil {
		return nil, err
	}
	x := &searchFindArtistsClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Search_FindArtistsClient interface {
	Recv() (*ArtistResponse, error)
	grpc.ClientStream
}

type searchFindArtistsClient struct {
	grpc.ClientStream
}

func (x *searchFindArtistsClient) Recv() (*ArtistResponse, error) {
	m := new(ArtistResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// SearchServer is the server API for Search service.
// All implementations must embed UnimplementedSearchServer
// for forward compatibility
type SearchServer interface {
	FindAlbums(*SearchMsg, Search_FindAlbumsServer) error
	FindTracks(*SearchMsg, Search_FindTracksServer) error
	FindPlaylists(*SearchMsg, Search_FindPlaylistsServer) error
	FindArtists(*SearchMsg, Search_FindArtistsServer) error
	mustEmbedUnimplementedSearchServer()
}

// UnimplementedSearchServer must be embedded to have forward compatible implementations.
type UnimplementedSearchServer struct {
}

func (UnimplementedSearchServer) FindAlbums(*SearchMsg, Search_FindAlbumsServer) error {
	return status.Errorf(codes.Unimplemented, "method FindAlbums not implemented")
}
func (UnimplementedSearchServer) FindTracks(*SearchMsg, Search_FindTracksServer) error {
	return status.Errorf(codes.Unimplemented, "method FindTracks not implemented")
}
func (UnimplementedSearchServer) FindPlaylists(*SearchMsg, Search_FindPlaylistsServer) error {
	return status.Errorf(codes.Unimplemented, "method FindPlaylists not implemented")
}
func (UnimplementedSearchServer) FindArtists(*SearchMsg, Search_FindArtistsServer) error {
	return status.Errorf(codes.Unimplemented, "method FindArtists not implemented")
}
func (UnimplementedSearchServer) mustEmbedUnimplementedSearchServer() {}

// UnsafeSearchServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to SearchServer will
// result in compilation errors.
type UnsafeSearchServer interface {
	mustEmbedUnimplementedSearchServer()
}

func RegisterSearchServer(s grpc.ServiceRegistrar, srv SearchServer) {
	s.RegisterService(&Search_ServiceDesc, srv)
}

func _Search_FindAlbums_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(SearchMsg)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(SearchServer).FindAlbums(m, &searchFindAlbumsServer{stream})
}

type Search_FindAlbumsServer interface {
	Send(*AlbumResponse) error
	grpc.ServerStream
}

type searchFindAlbumsServer struct {
	grpc.ServerStream
}

func (x *searchFindAlbumsServer) Send(m *AlbumResponse) error {
	return x.ServerStream.SendMsg(m)
}

func _Search_FindTracks_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(SearchMsg)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(SearchServer).FindTracks(m, &searchFindTracksServer{stream})
}

type Search_FindTracksServer interface {
	Send(*TrackResponse) error
	grpc.ServerStream
}

type searchFindTracksServer struct {
	grpc.ServerStream
}

func (x *searchFindTracksServer) Send(m *TrackResponse) error {
	return x.ServerStream.SendMsg(m)
}

func _Search_FindPlaylists_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(SearchMsg)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(SearchServer).FindPlaylists(m, &searchFindPlaylistsServer{stream})
}

type Search_FindPlaylistsServer interface {
	Send(*PlaylistResponse) error
	grpc.ServerStream
}

type searchFindPlaylistsServer struct {
	grpc.ServerStream
}

func (x *searchFindPlaylistsServer) Send(m *PlaylistResponse) error {
	return x.ServerStream.SendMsg(m)
}

func _Search_FindArtists_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(SearchMsg)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(SearchServer).FindArtists(m, &searchFindArtistsServer{stream})
}

type Search_FindArtistsServer interface {
	Send(*ArtistResponse) error
	grpc.ServerStream
}

type searchFindArtistsServer struct {
	grpc.ServerStream
}

func (x *searchFindArtistsServer) Send(m *ArtistResponse) error {
	return x.ServerStream.SendMsg(m)
}

// Search_ServiceDesc is the grpc.ServiceDesc for Search service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Search_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "auth.Search",
	HandlerType: (*SearchServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "FindAlbums",
			Handler:       _Search_FindAlbums_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "FindTracks",
			Handler:       _Search_FindTracks_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "FindPlaylists",
			Handler:       _Search_FindPlaylists_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "FindArtists",
			Handler:       _Search_FindArtists_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "search.proto",
}
