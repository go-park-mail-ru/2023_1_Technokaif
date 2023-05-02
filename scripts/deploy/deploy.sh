#!/bin/bash

docker login

docker push technokaif/fluire_api:latest
docker push technokaif/fluire_auth:latest
docker push technokaif/fluire_search:latest
docker push technokaif/fluire_user:latest
