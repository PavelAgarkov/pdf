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
    echo "fast_update_service - быстрый способ обновить сервис, если он уже работал" &&
    echo "docker_install - устанавливает docker и docker-compose"
    echo "ufw_init - запускает сетевой экран и устанавливает открытые порты"
    echo "monitor_ports - посмотреть открытые порт и сетевой экран"
}

#176.119.159.215 pdf-lifeguard.com www.pdf-lifeguard.com

generate_ssl() {
    apache2_stop &&
    apt install nginx &&
    nginx_enable &&
    apt-get install certbot &&
    apt-get install python3-certbot-nginx &&
    cd /etc/nginx/sites-available/ &&
    touch pdf-lifeguard.conf &&
    echo "server {
              listen 80;

              client_max_body_size 100m;
              server_tokens off;

              server_name pdf-lifeguard.com www.pdf-lifeguard.com;

              location / {
                  proxy_connect_timeout 200s;
                  proxy_send_timeout 200s;
                  proxy_read_timeout 200s;

                  proxy_set_header X-NginX-Proxy true;
                  client_max_body_size 100m;

                  proxy_pass http://0.0.0.0:3000/;
              }
          }" >> /etc/nginx/sites-available/pdf-lifeguard.conf &&
    ln -s /etc/nginx/sites-available/pdf-lifeguard.conf /etc/nginx/sites-enabled/ &&
    certbot --nginx -d pdf-lifeguard.com -d www.pdf-lifeguard.com &&
    nginx -t && systemctl restart nginx
}

regenerate_ssl() {
    certbot --nginx -d pdf-lifeguard.com -d www.pdf-lifeguard.com &&
    nginx -t &&
    systemctl restart nginx
}

monitor_ports() {
  nmap -4 -Pn 176.119.159.215 &&
  ufw status numbered &&
  netstat -tulpn | grep LISTEN
}

ssh_gen() {
  cd ~/ &&
  mkdir  ~/.ssh && cd ~/.ssh &&
  ssh-keygen -t ed25519 -C "agarkov.ru@mail.ru" &&
  eval "$(ssh-agent -s)" && ssh-add ~/.ssh/pdf
}

ssh_update() {
     cd /var/www/pdf && eval "$(ssh-agent -s)" && ssh-add ~/.ssh/pdf
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

git_install() {
  apt-get update &&
  apt install git
}

apache2_stop() {
  /etc/init.d/apache2 stop &&
  update-rc.d apache2 disable &&
  echo "apache2 stopped complete"

  systemctl stop mysql.service &&
  systemctl disable mysql &&
  apt autoremove && apt autoclean &&
  rm -rf /etc/mysql /var/lib/mysql &&
  apt-get purge mysql-server &&
  echo "mysql stopped complete"
}

nginx_enable() {
  /etc/init.d/nginx start &&
  update-rc.d nginx enable
}

start_service() {
  apache2_stop &&
  nginx_enable &&
  git_update &&
  stop_service &&
  docker-compose -f /var/www/pdf/docker-compose-prode.yaml build --remove-orphans &&
  build_project &&
  docker-compose -f /var/www/pdf/docker-compose-prode.yaml up -d --remove-orphans &&
  echo "service start on port 443" &&
  docker-compose -f /var/www/pdf/docker-compose-prode.yaml stop node
}

fast_update_service() {
  git_update &&
  build_project &&
  stop_service &&
  docker-compose -f /var/www/pdf/docker-compose-prode.yaml up -d  --remove-orphans &&
  echo "service start on port 443" &&
  docker-compose -f /var/www/pdf/docker-compose-prode.yaml stop node
}

stop_service() {
  docker-compose -f /var/www/pdf/docker-compose-prode.yaml stop &&
  echo "service stopped complete"
}

start_after_stop() {
    docker-compose -f /var/www/pdf/docker-compose-prode.yaml up -d  --remove-orphans &&
    echo "service start on port 443" &&
    docker-compose -f /var/www/pdf/docker-compose-prode.yaml stop node
}

docker_install() {
  apt install docker.io &&
  apt install docker-compose &&
  apt install net-tools &&
  apt install htop
}

git_init() {
  git_install &&
  cd /var/www/ &&
  git clone git@github.com:PavelAgarkov/pdf.git &&
  cd pdf/ &&
  git clone git@github.com:PavelAgarkov/pdf-frontend.git &&
  chmod 667 /var/www/pdf/pdf-frontend &&
  chmod 666 /var/www/pdf/pdf-frontend/package-lock.json &&
  chmod 666 /var/www/pdf/pdf-frontend/package.json
}

ufw_init() {
  apt install nmap &&
  apt-get install ufw &&
  ufw enable &&
  ufw allow https &&
  ufw allow http &&
  ufw allow ssh &&
  ufw allow 'Nginx Full' &&
  ufw deny 3000/tcp &&
  ufw status numbered
}

go_install() {
  apt install curl &&
  cd /home &&
  curl -OL https://golang.org/dl/go1.21.3.linux-amd64.tar.gz &&
  tar -C /usr/local -xvf go1.21.3.linux-amd64.tar.gz &&
  cd /usr/local/go &&
  echo "export PATH=$PATH:/usr/local/go/bin" >> ~/.profile &&
  cd /var/www/pdf &&
  source ~/.profile &&
  backend_build
}

if declare -f "$1" > /dev/null
then
  "$@"
else
  echo "'$1' is not a known function name" >&2
  exit
fi