#!/bin/bash

fuser -k 8082/tcp
docker-compose up -d
docker-compose run go-api

