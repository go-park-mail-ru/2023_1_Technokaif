#!/bin/bash

#TAG="$(git branch --show-current)_$(git rev-parse --short HEAD)"
if [[ ! $TAG ]]; then
    echo "Using TAG=latest"
    TAG=latest
fi

docker build -t fluire_search:$TAG -f Dockerfile.application --build-arg APP=search/search.go --build-arg PORT=4442 . 
docker tag fluire_search:$TAG technokaif/fluire_search:$TAG
