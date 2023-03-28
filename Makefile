.PHONY: all

all: clear_media server_start

server_start:
	go run ./cmd/app/main.go

clear_media:
	rm -r ./img
