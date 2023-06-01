#!/bin/bash

#TAG="$(git branch --show-current)_$(git rev-parse --short HEAD)"
if [[ ! $TAG ]]; then
    echo "Using TAG=latest"
    TAG=latest
fi

docker push technokaif/fluire_api:$TAG
docker push technokaif/fluire_auth:$TAG
docker push technokaif/fluire_search:$TAG
docker push technokaif/fluire_user:$TAG
