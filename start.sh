#!/usr/bin/env sh
export USERID=$(id -u)
docker-compose up -d
gulp
