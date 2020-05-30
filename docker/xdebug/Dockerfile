FROM mygento/php:7.1-full

ARG HOST_ARG

RUN echo "\033[37;1;42m \033[5mCheck your host. Current host is $HOST_ARG \033[0m"

RUN yes | pecl install xdebug-2.6.1 && \
    echo zend_extension=/usr/lib/php/20160303/xdebug.so >> /etc/php/7.1/mods-available/xdebug.ini && phpenmod xdebug && \
    echo xdebug.remote_port=9001 >> /etc/php/7.1/mods-available/xdebug.ini && \
    echo xdebug.remote_connect_back=1 >> /etc/php/7.1/mods-available/xdebug.ini && \
    echo xdebug.remote_host=$HOST_ARG >> /etc/php/7.1/mods-available/xdebug.ini && \
    echo xdebug.remote_enable=1 >> /etc/php/7.1/mods-available/xdebug.ini && \
    echo xdebug.idekey=netbeans-xdebug >> /etc/php/7.1/mods-available/xdebug.ini && \
    echo xdebug.show_local_vars=0 >> /etc/php/7.1/mods-available/xdebug.ini && \
    echo xdebug.output_buffering=off >> /etc/php/7.1/mods-available/xdebug.ini && \
    echo xdebug.cli_color=1 >> /etc/php/7.1/mods-available/xdebug.ini && \
    echo xdebug.remote_log=/var/log/xdebug.log >> /etc/php/7.1/mods-available/xdebug.ini
