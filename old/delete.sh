#!/usr/bin/env sh
export USERID=$(id -u)
docker-compose stop
docker-compose rm
