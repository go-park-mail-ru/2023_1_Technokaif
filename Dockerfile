FROM golang:1.20

WORKDIR /app

COPY . .

RUN go build cmd/api/main.go

EXPOSE 4444

CMD ["./main"]