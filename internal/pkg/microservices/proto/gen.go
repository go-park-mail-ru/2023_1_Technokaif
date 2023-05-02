package proto

//go:generate protoc --proto_path=. --go_out=.. --go-grpc_out=.. auth.proto
//go:generate protoc --proto_path=. --go_out=.. --go-grpc_out=.. user.proto
//go:generate protoc --proto_path=. --go_out=.. --go-grpc_out=.. search.proto
//go:generate protoc --proto_path=. --go_out=.. --go-grpc_out=.. --go_opt=module=github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/microservices common.proto
