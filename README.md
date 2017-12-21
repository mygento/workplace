# Magento Local worker image


PHP [![](https://images.microbadger.com/badges/image/mygento/php:7.1-fpm.svg)](https://microbadger.com/images/mygento/php:7.1-fpm)

NGINX [![](https://images.microbadger.com/badges/image/luckyraul/nginx:backports.svg)](https://microbadger.com/images/luckyraul/nginx:backports)

MYSQL
[![](https://images.microbadger.com/badges/image/mygento/mysql:5.6.svg)](https://microbadger.com/images/mygento/mysql:5.6)

DB HOST: db

DB Name: magento

DB user/pass: mygento

RUN magento bin:

docker exec -ti -u www-data m2-php /var/www/public/bin/magento
