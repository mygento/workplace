#!/usr/bin/env sh
export USERID=$(id -u)
sudo usermod -a -G www-data ${USER}
rm -f application/composer.lock
composer install -d application --ignore-platform-reqs
npm install
docker-compose pull
