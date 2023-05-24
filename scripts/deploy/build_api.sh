#!/bin/bash

#TAG="$(git branch --show-current)_$(git rev-parse --short HEAD)"
if [[ ! $TAG ]]; then
    echo "Using TAG=latest"
    TAG=latest
fi

docker build -t fluire_api:$TAG -f Dockerfile.application --build-arg APP=api/main.go --build-arg PORT=4444 .
docker tag fluire_api:$TAG technokaif/fluire_api:$TAG
