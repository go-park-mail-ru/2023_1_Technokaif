#!/bin/bash

#TAG="$(git branch --show-current)_$(git rev-parse --short HEAD)"
#$TAG=latest

docker push technokaif/fluire_api:$TAG
docker push technokaif/fluire_auth:$TAG
docker push technokaif/fluire_search:$TAG
docker push technokaif/fluire_user:$TAG
