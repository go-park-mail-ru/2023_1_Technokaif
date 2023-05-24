.PHONY: all

# all: clear_media server_start

build:
	./scripts/deploy/build_all.sh

push:
	./scripts/deploy/push.sh

start:
	docker-compose down
	docker-compose up -d
	docker-compose up

clean_containers:
	docker system prune

clear_media:
	rm -r ./img ./covers ./records ./avatars

lint:
	go vet ./... && golangci-lint run

check_coverage:
	go test -coverpkg=./... -coverprofile=coverage.out ./... \
	&& cat coverage.out | fgrep -v "mocks" | fgrep -v "docs" > purified_coverage.out \
	&& go tool cover -func purified_coverage.out | grep total

check_html_coverage:
	go test -coverpkg=./... -coverprofile=coverage.out ./... \
	&& cat coverage.out | fgrep -v "mocks" | fgrep -v "docs" > purified_coverage.out \
	&& go tool cover -func purified_coverage.out | grep total \
	&& go tool cover -html purified_coverage.out -o cover.html 

generate_api_docs:
	swag init -g cmd/api/main.go

generate:
	go generate ./...
