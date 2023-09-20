#!/bin/zsh -e

using() {
    echo "Укажите команду при запуске: ./monolit.sh [command]"
}

go_build() {
  go mod tidy &&
  go mod vendor &&
  go build
}

go_run_after_build() {
  go_build &&
  ./pdf
}

frontend_build() {
  cd frontend/ &&
  pwd &&
  nvm use 18 &&
  npm run build &&
  cd ..
}


operation="$1"
if [[ -z "$command" ]]; then
    using
    exit 0
else
    # shellcheck disable=SC2068
    $command $@
fi