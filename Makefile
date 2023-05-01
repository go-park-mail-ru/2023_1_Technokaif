.PHONY: all

# all: clear_media server_start

build:
	docker-compose down
	docker-compose build

start:
	docker-compose down
	docker-compose up -d
	docker-compose up

drop_db:
	rm -r ./.pgdata

clean_containers:
	docker system prune

api_start:
	go run ./cmd/api/main.go

auth_start:
	go run ./cmd/auth/auth.go

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
	&& go tool cover -html=purified_coverage.out 

generate_api_docs:
	swag init -g cmd/api/main.go

generate:
	go generate ./...
