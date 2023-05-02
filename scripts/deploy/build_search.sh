#!/bin/bash

docker build -t fluire_search -f Dockerfile.application --build-arg APP=search/search.go --build-arg PORT=4442 . 
docker tag fluire_search:latest technokaif/fluire_search:latest
