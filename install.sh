#!/usr/bin/env sh
export USERID=$(id -u)
rm -f application/composer.lock
composer install -d application --ignore-platform-reqs
npm install
docker-compose pull
