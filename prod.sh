#!/bin/bash -e

#anysync() {
#    apt-get update
#      apt install git
#  cd ~/.ssh && ssh-keygen
#  ssh-keygen -t rsa
#   eval "$(ssh-agent -s)"
#   ssh-add ~/.ssh/id_ed25519
#   git remote set-url origin git@github.com:PavelAgarkov/pdf-frontend.git


#  /etc/init.d/apache2 stop
#  update-rc.d apache2 disable

#  apt install docker.io
#  apt install docker-compose
#  apt install net-tools


#  cd /home
#  curl -OL https://golang.org/dl/go1.21.3.linux-amd64.tar.gz
#  cd /usr/local/go

#  nano ~/.profile
  #export PATH=$PATH:/usr/local/go/bin
#  cd /var/www/pdf
#  source ~/.profile

#  go build && go mod vendor && go mod tidy
#}

git_init() {
        apt install git &&
           cd ~/.ssh
#           ssh-keygen -t rsa &&
#            eval "$(ssh-agent -s)" && ssh-add ~/.ssh/id_ed25519
#  cd /var/www/ && git clone git@github.com:PavelAgarkov/pdf.git && cd pdf/ && git clone git@github.com:PavelAgarkov/pdf-frontend.git

#chmod -r 666 /var/www/pdf/pdf-frontend
}

git_update() {
  cd /var/www/pdf &&
  git pull &&
  cd /pdf-frontend/ &&
  git pull
}

build_project() {
  backend_build && frontend_build
}

backend_build() {
  cd /var/www/pdf &&
       go build &&
       go mod vendor &&
       go mod tidy
}

frontend_build() {
   docker exec node-local npm install &&
  docker exec node-local npm run build
}

go_install() {
    cd /home &&
    curl -OL https://golang.org/dl/go1.21.3.linux-amd64.tar.gz &&
    cd /usr/local/go &&
    tar -C /usr/local -xvf go1.21.3.linux-amd64.tar.gz &&
    echo "export PATH=$PATH:/usr/local/go/bin" >> ~/.profile &&
    cd /var/www/pdf &&
    source ~/.profile &&
    backend_build
}

git_install() {
      apt-get update &&
        apt install git
}

apache2_stop() {
    /etc/init.d/apache2 stop &&
    update-rc.d apache2 disable
}

start_service() {
 stop_service &&
  docker-compose -f /var/www/pdf/docker-compose-prode.yaml build &&
  docker-compose -f /var/www/pdf/docker-compose-prode.yaml up -d
}

stop_service() {
  docker-compose -f /var/www/pdf/docker-compose-prode.yaml stop
}

docker_install() {
    apt install docker.io &&
    apt install docker-compose &&
    apt install net-tools
}

if declare -f "$1" > /dev/null
then
  "$@"
else
  echo "'$1' is not a known function name" >&2
  exit
fi