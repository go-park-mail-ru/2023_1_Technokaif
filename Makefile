.PHONY: all

all: clear_media server_start

api_start:
	go run ./cmd/api/main.go

auth_start:
	go run ./cmd/auth/auth.go

clear_media:
	rm -r ./img

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
	swag init -g cmd/app/main.go

generate_mocks:
	go generate ./...
