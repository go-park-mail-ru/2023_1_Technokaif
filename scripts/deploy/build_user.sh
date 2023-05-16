#!/bin/bash

#TAG="$(git branch --show-current)_$(git rev-parse --short HEAD)"
#TAG=latest

docker build -t fluire_user:$TAG -f Dockerfile.application --build-arg APP=user/user.go --build-arg PORT=4441 . 
docker tag fluire_user:$TAG technokaif/fluire_user:$TAG
