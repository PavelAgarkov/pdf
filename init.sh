#!/bin/bash -e


apt install git

apt-get update

apt install docker.io
apt install docker-compose
apt install net-tools
cd /var/www/ && git clone https://github.com/PavelAgarkov/pdf.git && cd pdf/ && git clone