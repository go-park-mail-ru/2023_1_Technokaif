#!/bin/bash

#TAG="$(git branch --show-current)_$(git rev-parse --short HEAD)"
if [[ ! $TAG ]]; then
    echo "Using TAG=latest"
    TAG=latest
fi

docker build -t fluire_auth:$TAG -f Dockerfile.application --build-arg APP=auth/auth.go --build-arg PORT=4443 . 
docker tag fluire_auth:$TAG technokaif/fluire_auth:$TAG
