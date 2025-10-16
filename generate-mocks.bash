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

mockgen -destination=internal/mocks/mock_user_reciever.go -package=mocks ./internal/handlers UserReceiver
mockgen -destination=internal/mocks/mock_user_registerer.go -package=mocks ./internal/handlers UserRegisterer
mockgen -destination=internal/mocks/mock_user_loginer.go -package=mocks ./internal/handlers UserLoginer
mockgen -destination=internal/mocks/mock_token_parser.go -package=mocks ./internal/middlewares TokenParser

mockgen -destination=internal/mocks/mock_orders_service.go -package=mocks ./internal/handlers OrdersService

mockgen -destination=internal/mocks/mock_accrualer.go -package=mocks ./internal/services Accrualer
mockgen -destination=internal/mocks/mock_withdrawer.go -package=mocks ./internal/handlers Withdrawer
mockgen -destination=internal/mocks/mock_balancer.go -package=mocks ./internal/handlers Balancer

mockgen -destination=internal/mocks/mock_order_job_processor.go -package=mocks ./internal/services OrderJobProcessor
mockgen -destination=internal/mocks/mock_orders_queue.go -package=mocks ./internal/handlers OrdersQueue

echo "Finish"
