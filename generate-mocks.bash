#!/usr/bin/env bash

# Прерывать скрипт при ошибке
set -e

# Прерывать если не передана переменная
set -u

CURRENT_DIR_BASENAME=`basename $PWD`

# Скрипт нужно запускать из корня проекта, чтобы не возиться с путями в командах.
# Предполагаю что репа названа как на гитхабе. Учитывать ренейминг лень.
if [[ $CURRENT_DIR_BASENAME != "gophermart" ]]; then
  echo "Run only from root project dir"
  exit 1
fi

mockgen -destination=internal/mocks/mock_users_storage.go -package=mocks -mock_names Storage=MockUsersStorage ./internal/storage/users Storage

mockgen -destination=internal/mocks/mock_services_auther.go -package=mocks ./internal/services Auther

echo "Finish"
