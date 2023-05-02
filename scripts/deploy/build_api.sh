#!/bin/bash

docker build -t fluire_api -f Dockerfile.application --build-arg APP=api/main.go --build-arg PORT=4444 .
docker tag fluire_api:latest technokaif/fluire_api:latest
