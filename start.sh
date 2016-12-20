#!/usr/bin/env sh
docker-compose up -d
chgrp -R www-data application
# chmod a+w application/public/app/etc
# chmod -R a+w application/public/media
gulp
