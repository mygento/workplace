#!/usr/bin/env sh
sudo groupadd -g 33 www-data
sudo useradd --uid 33 --gid 33 -d /var/www/ -s /bin/sh www-data
sudo usermod -a -G www-data ${USER}
rm -f application/composer.lock
composer install -d application --ignore-platform-reqs
npm install
docker-compose pull
