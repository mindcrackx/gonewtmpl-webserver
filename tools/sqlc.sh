#! /usr/bin/env bash

docker run -v ./sqlc.yaml:/app/sqlc.yaml -v ./internal/db:/app/internal/db $(docker build -q tools/sqlc/)
