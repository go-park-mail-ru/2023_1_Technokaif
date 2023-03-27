.PHONY: all

all: serverstart

serverstart:
	go run ./cmd/app/main.go
