#!/usr/bin/env sh
docker-compose up -d
chmod g+s application
chmod -R g+w application/public/media
chmod -R g+w application/public/app/etc
chmod -R g+w application/public/var
gulp
