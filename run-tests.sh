#!/usr/bin/env bash

CONTAINER_ID=$(docker run --rm --detach -p "27017:27017" mongo:latest)
MONGO_URL="mongodb://localhost:27017/hf?ssl=false"
DEPLOY_ENV=testing
PORT=8080

export CONTAINER_ID MONGO_URL DEPLOY_ENV PORT

go test ./...

docker container kill "${CONTAINER_ID}"
