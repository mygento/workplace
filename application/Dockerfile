FROM busybox
RUN mkdir -p /var/www
WORKDIR /var/www
ADD ./public /var/www
RUN chown 33:33 -R /var/www
CMD ['/var/www']