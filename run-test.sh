#!/bin/bash

docker-compose up -d
docker-compose run --entrypoint="go test" go-api

