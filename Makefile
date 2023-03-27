.PHONY: all

all: dbstop dbstart serverstart

serverstart:
	go run ./cmd/app/main.go

dbstart:
	doocker-compose up

dbstop:
	doocker-compose down