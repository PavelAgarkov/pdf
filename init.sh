#!/bin/bash -e


apt install git

apt-get update

apt install docker.io
apt install docker-compose
apt install net-tools
cd /var/www/ && git clone https://github.com/PavelAgarkov/pdf.git && cd pdf/ && git clone https://github.com/PavelAgarkov/pdf-frontend.git


cd /home
curl -OL https://golang.org/dl/go1.21.3.linux-amd64.tar.gz
cd /usr/local/go

nano ~/.profile
#export PATH=$PATH:/usr/local/go/bin
cd /var/www/pdf
source ~/.profile

go build && go mod vendor && go mod tidy