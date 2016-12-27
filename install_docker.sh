#!/usr/bin/env sh
sudo apt-get update && sudo apt-get install apt-transport-https ca-certificates gnupg2
sudo apt-key adv \
       --keyserver hkp://ha.pool.sks-keyservers.net:80 \
       --recv-keys 58118E89F3A912897C070ADBF76221572C52609D
sudo echo 'deb https://apt.dockerproject.org/repo debian-jessie main' > /etc/apt/sources.list.d/docker.list
sudo apt-get update && sudo apt-get install docker-engine
sudo groupadd docker
sudo gpasswd -a ${USER} docker
sudo service docker restart
wget "https://github.com/docker/compose/releases/download/1.9.0/docker-compose-$(uname -s)-$(uname -m)" -O docker-compose
sudo mv docker-compose /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose
