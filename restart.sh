#!/bin/bash

env GOOS=linux GOARCH=amd64 gb build
docker-compose build
docker-compose up -d
