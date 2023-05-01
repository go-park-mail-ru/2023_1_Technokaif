FROM golang:1.20-alpine AS build_stage

LABEL maintainer="yarik1448kuzmin@gmail.com"

WORKDIR /app
COPY . .
RUN go build -o ./out/auth_bin cmd/auth/auth.go
RUN go build -o ./out/api_bin cmd/api/main.go

# 2 шаг
FROM alpine AS run_stage

WORKDIR /out
COPY --from=build_stage /app/out/auth_bin /out/
RUN chmod +x ./auth_bin
EXPOSE 4443/tcp

COPY --from=build_stage /app/out/api_bin /out/
RUN chmod +x ./api_bin
EXPOSE 4444/tcp
