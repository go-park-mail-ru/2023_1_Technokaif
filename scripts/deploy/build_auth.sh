#!/bin/bash

docker build -t fluire_auth -f Dockerfile.application --build-arg APP=auth/auth.go --build-arg PORT=4443 . 
docker tag fluire_auth:latest yarik_tri/fluire_auth:latest
