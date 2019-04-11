#!/usr/bin/env sh
sudo apt-get -qq update && sudo apt-get install -qqy wget lsb-release apt-transport-https ca-certificates gnupg2 curl
sudo apt-get remove -y docker docker-engine docker.io
DISTR="$(lsb_release -is)"
if [ "$DISTR" == "Debian" ]; then
  curl -fsSL https://download.docker.com/linux/debian/gpg | sudo apt-key add -
  echo "deb [arch=amd64] https://download.docker.com/linux/debian $(lsb_release -cs) stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null
fi
if [ "$DISTR" == "Ubuntu" ]; then
  curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -
  echo "deb [arch=amd64] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null
fi
sudo apt-get -qq update && sudo apt-get install -qqy docker-ce
sudo groupadd docker
sudo usermod -aG docker $USER
sudo service docker restart
wget "https://github.com/docker/compose/releases/download/1.23.2/docker-compose-$(uname -s)-$(uname -m)" -O docker-compose
sudo mv docker-compose /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose
sudo usermod -aG www-data $USER
