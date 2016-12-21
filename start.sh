#!/usr/bin/env sh
docker-compose up -d
chgrp -R www-data application
chmod -R a+w application/app/etc
chmod -R a+w application/var
chmod -R a+w application/pub/media
chmod -R a+w application/pub/static
gulp
