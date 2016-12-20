#!/usr/bin/env sh
sudo groupadd -g 33 www-data
sudo useradd --uid 33 --gid 33 -d /var/www/ -s /bin/sh www-data
sudo usermod -a -G www-data ${USER}
composer install
npm install
docker-compose build
docker-compose up -d
chgrp -R www-data application
gulp
