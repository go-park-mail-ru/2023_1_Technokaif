FROM golang:1.20-alpine

# Set the Current Working Directory inside the container
WORKDIR /api

# We want to populate the module cache based on the go.{mod,sum} files.
COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

# Build the Go app
RUN go build -o ./out/api ./cmd/api/main.go


# This container exposes port 8080 to the outside world
EXPOSE 4444

# Run the binary program produced by `go install`
CMD ["./out/api"]