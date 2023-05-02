#!/bin/bash

docker build -t fluire_user -f Dockerfile.application --build-arg APP=user/user.go --build-arg PORT=4441 . 
docker tag fluire_user:latest technokaif/fluire_user:latest
