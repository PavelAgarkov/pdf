#!/bin/bash -e

command() {
    echo "Укажите команду при запуске: ./help.sh [command]" &&
    echo "git_init - инициализирует гит и проект, выкачивает и устанавливает права на фронтенд" &&
    echo "ssh_gen - геренерирует пару ssh и добавляет в агента"  &&
    echo "git_update - обновляет изменения в мастер ветках бэка и фронта" &&
    echo "build_project - выполняет сборку фронтенда и бэкенда" &&
    echo "backend_build - сборка только бэкенда" &&
    echo "frontend_build - сборка только фронтенда" &&
    echo "go_install - выкачивает го 1.21.3 и распаковывает его с добавлением в пас" &&
    echo "git_install - скачивает гит" &&
    echo "apache2_stop - останавливает апач2 и выводит из автозагрузки, чтобы не подавлял nginx в контейнере" &&
    echo "start_service - запускает сервис со всеми обновления и пересборками" &&
    echo "stop_service - останавливает все контенеры" &&
    echo "docker_install - устанавливает docker и docker-compose"
}

git_init() {
  apt install git &&
  cd /var/www/ &&
  git clone git@github.com:PavelAgarkov/pdf.git &&
  cd pdf/ &&
  git clone git@github.com:PavelAgarkov/pdf-frontend.git &&
  chmod 666 /var/www/pdf/pdf-frontend
}

ssh_gen() {
  cd ~/.ssh &&
  ssh-keygen -t rsa &&
   eval "$(ssh-agent -s)"
#   ssh-add ~/.ssh/id_ed25519
}

git_update() {
  cd /var/www/pdf &&
  git pull &&
  echo "backend pull complete" &&
  cd /var/www/pdf/pdf-frontend/ &&
  git pull &&
  echo "frontend pull complete"
}

build_project() {
  frontend_build && backend_build
}

backend_build() {
  cd /var/www/pdf &&
       go build &&
       go mod vendor &&
       go mod tidy &&
       echo "backend build complete"
}

frontend_build() {
  docker-compose -f /var/www/pdf/docker-compose-prode.yaml start node &&
   docker exec node-local npm install &&
  docker exec node-local npm run build &&
  docker-compose -f /var/www/pdf/docker-compose-prode.yaml stop node &&
  echo "frontend build complete"
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
    update-rc.d apache2 disable &&
    echo "apache2 stopped complete"
}

start_service() {
  apache2_stop &&
  git_update &&
  stop_service &&
  docker-compose -f /var/www/pdf/docker-compose-prode.yaml build &&
  build_project &&
  docker-compose -f /var/www/pdf/docker-compose-prode.yaml up -d &&
  echo "service start on port 80" &&
  docker-compose -f /var/www/pdf/docker-compose-prode.yaml stop node
}

fast_update_service() {
  git_update &&
  build_project &&
  stop_service &&
  docker-compose -f /var/www/pdf/docker-compose-prode.yaml up -d &&
  echo "service start on port 80" &&
  docker-compose -f /var/www/pdf/docker-compose-prode.yaml stop node
}

stop_service() {
  docker-compose -f /var/www/pdf/docker-compose-prode.yaml stop &&
  echo "service stopped complete"
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