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
mockgen -destination=internal/mocks/mock_orders_storage.go -package=mocks -mock_names Storage=MockOrdersStorage ./internal/storage/orders Storage
mockgen -destination=internal/mocks/mock_balance_storage.go -package=mocks -mock_names Storage=MockBalanceStorage ./internal/storage/balance Storage

mockgen -destination=internal/mocks/mock_services_auther.go -package=mocks ./internal/services Auther

mockgen -destination=internal/mocks/mock_services_accrual_service.go -package=mocks ./internal/services IAccrualService
mockgen -destination=internal/mocks/mock_services_orders_service.go -package=mocks ./internal/services IOrdersService
mockgen -destination=internal/mocks/mock_services_balancer.go -package=mocks ./internal/services IBalancer

mockgen -destination=internal/mocks/mock_services_orders_queue.go -package=mocks ./internal/services OrdersQueue

echo "Finish"
