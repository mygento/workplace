#!/usr/bin/env sh
export USERID=$(id -u)
if [ -f ./src/auth.json ]; then
    cp ./src/auth.json application
fi
rm -f application/composer.lock
composer install -d application --ignore-platform-reqs
npm install
docker-compose pull
