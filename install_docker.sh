#!/usr/bin/env sh
sudo apt-get -qq update && sudo apt-get install -qqy wget apt-transport-https ca-certificates gnupg2 curl
curl -fsSL https://download.docker.com/linux/debian/gpg | sudo apt-key add -
sudo apt-get remove -y docker docker-engine
sudo sh -c "echo deb [arch=amd64] https://download.docker.com/linux/debian jessie stable > /etc/apt/sources.list.d/docker.list"
sudo apt-get -qq update && sudo apt-get install -qqy docker-ce
sudo groupadd docker
sudo gpasswd -a ${USER} docker
sudo service docker restart
wget "https://github.com/docker/compose/releases/download/1.18.0/docker-compose-$(uname -s)-$(uname -m)" -O docker-compose
sudo mv docker-compose /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose
sudo usermod -a -G www-data ${USER}
