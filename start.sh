#!/usr/bin/env sh
composer install
npm install
docker-compose build
