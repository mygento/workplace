# Magento Local worker image


PHP [![](https://images.microbadger.com/badges/image/mygento/php:5.6-fpm.svg)](https://microbadger.com/images/mygento/php:5.6-fpm)

NGINX [![](https://images.microbadger.com/badges/image/mygento/nginx:backports.svg)](https://microbadger.com/images/mygento/nginx:backports)

MYSQL
[![](https://images.microbadger.com/badges/image/mygento/mysql:5.6.svg)](https://microbadger.com/images/mygento/mysql:5.6)


REDIS
[![](https://images.microbadger.com/badges/image/redis:3.svg)](https://microbadger.com/images/redis:3)


DB HOST: db

DB Name: magento

DB user/pass: mygento

RUN magento bin:

docker exec -ti m2-php /var/www/public/bin/magento

CHOWN (before start):

sudo chown -R hamster:www-data application/
