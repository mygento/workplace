#!/usr/bin/env sh
export USERID=$(id -u)
docker-compose up -d
docker exec -ti m2-php chmod +x /var/www/public/bin/magento
gulp
