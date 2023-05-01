FROM golang:1.20-alpine AS build_stage

LABEL maintainer="yarik1448kuzmin@gmail.com"

WORKDIR /app
COPY . .
RUN go build -o ./out/auth_bin cmd/auth/auth.go

# 2 шаг
FROM alpine AS run_stage

WORKDIR /out
COPY --from=build_stage /app/out/auth_bin /out/
RUN chmod +x ./auth_bin
EXPOSE 4443/tcp

ENTRYPOINT ./auth_bin
