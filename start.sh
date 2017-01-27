#!/usr/bin/env sh
docker-compose up -d
docker exec -ti m2-php chmod +x /var/www/public/bin/magento
chgrp -R www-data application
chmod g+s application
chmod -R a+w application/app/etc
chmod -R a+w application/var
chmod -R a+w application/pub/media
chmod -R a+w application/pub/static
gulp
