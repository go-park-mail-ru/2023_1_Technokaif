package proto

//go:generate protoc --go_out=./generated --go-grpc_out=./generated --go-grpc_opt=paths=source_relative --go_opt=paths=source_relative auth.proto
