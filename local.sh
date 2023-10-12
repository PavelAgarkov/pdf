#!/bin/bash -e

. ~/.nvm/nvm.sh
. ~/.profile
. ~/.bashrc

using() {
    echo "Укажите команду при запуске: ./help.sh [command]"
}

go_build() {
  go mod tidy &&
  go mod vendor &&
  go build .
}

nwm_switch_to_18() {
  nvm use 18
}

start_pdf_with_build() {
  cd pdf-frontend/ &&
   pwd &&
    nwm_switch_to_18 &&
     npm run build &&
      cd ../ &&
       go build -race &&
        ./pdf
}


if declare -f "$1" > /dev/null
then
  "$@"
else
  echo "'$1' is not a known function name" >&2
  exit
fi