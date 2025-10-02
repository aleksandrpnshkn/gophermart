# Гофермарт

## Обновление шаблона
Для обновления кода автотестов:
```bash
git remote add -m master template https://github.com/yandex-praktikum/go-musthave-diploma-tpl.git
git fetch template && git checkout template/master .github
```

## БД
Окружение:
```bash
# Установить клиент для работы с БД (psql)
apt install postgresql-client

docker compose up --detach

# С хоста
psql --host 127.0.0.1 --port 5433 --username admin --password --dbname gophermart

# Для очистки базы
docker compose down --volumes
```

Для работы с миграциями установить migrate - https://github.com/golang-migrate/migrate/tree/v4.18.3/cmd/migrate . Затем в корне проекта:
```bash
~/golang-migrate/migrate create -ext sql -dir ./internal/storage/migrations -seq create_example_table

~/golang-migrate/migrate -database "postgres://admin:qwerty@localhost:5433/gophermart?sslmode=disable" -path ./internal/storage/migrations up
~/golang-migrate/migrate -database "postgres://admin:qwerty@localhost:5433/gophermart?sslmode=disable" -path ./internal/storage/migrations down
```

## Тестирование
Сгенерировать моки:
```bash
# из корня проекта
./generate-mocks.bash
```

Запустить unit-тесты (`-count` для отключения кэша, помогает отлавливать flaky-тесты):
```bash
go test -count=10 ./...
```

## Сборка и запуск
```bash
# Добавить go в PATH
source ~/.profile

# Запустить сервер
go build -o cmd/gophermart/gophermart cmd/gophermart/*go \
    && ./cmd/gophermart/gophermart -l debug
```

## Тестовые запросы
```bash
curl --include localhost:8081/
curl --include localhost:8081/api/ping

# регистрация
curl --request POST \
    --header "Content-Type: application/json" \
    --data '{"login": "user", "password": "secret"}' \
    --include \
    localhost:8081/api/user/register

# логин
curl --request POST \
    --header "Content-Type: application/json" \
    --data '{"login": "user", "password": "secret"}' \
    --include \
    localhost:8081/api/user/login

# добавить заказ в обработку заказ
curl --request POST \
    --header "Content-Type: text/plain" \
    --cookie "auth_token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySUQiOjF9.pL0mBx3adBKYmOpgyObgX1jl2XoifJnpFhiKKs4wgO0" \
    --data '12345678903' \
    --include \
    localhost:8081/api/user/orders 

# проверить баланс
curl --request GET \
    --header "Content-Type: application/json" \
    --cookie "auth_token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySUQiOjF9.pL0mBx3adBKYmOpgyObgX1jl2XoifJnpFhiKKs4wgO0" \
    --include \
    localhost:8081/api/user/balance

# оплатить заказ бонусами
curl --request POST \
    --header "Content-Type: application/json" \
    --cookie "auth_token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySUQiOjF9.pL0mBx3adBKYmOpgyObgX1jl2XoifJnpFhiKKs4wgO0" \
    --data '{"order": "12345678903", "sum": 123}' \
    --include \
    localhost:8081/api/user/balance/withdraw

# проверить списания
curl --request GET \
    --header "Content-Type: application/json" \
    --cookie "auth_token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySUQiOjF9.pL0mBx3adBKYmOpgyObgX1jl2XoifJnpFhiKKs4wgO0" \
    --include \
    localhost:8081/api/user/withdrawals
```

```sql
# принудительно накинуть бонусов для тестов
UPDATE orders 
SET accrual = 1000
WHERE number = '12345678903';

INSERT INTO balance_logs (id, order_number, user_id, amount, processed_at) 
VALUES (DEFAULT, '12345678903', 1, 1000, DEFAULT);
```

## accrual
```bash
# Запустить сервис расчёта баллов accrual
./cmd/accrual/accrual_linux_amd64 \
    -a "localhost:8083" \
    -d "postgres://admin:qwerty@localhost:5434/accrual?sslmode=disable"

# создать товар
curl --request POST \
    --header "Content-Type: application/json" \
    --data '{"match": "Bork", "reward": 10, "reward_type": "%"}' \
    --include \
    localhost:8083/api/goods

# зарегать заказ
curl --request POST \
    --header "Content-Type: application/json" \
    --data '{"order": "12345678903", "goods": [{"description": "Чай", "price": 7000}]}' \
    --include \
    localhost:8083/api/orders

# проверить статус заказа
curl --request GET \
    --include \
    localhost:8083/api/orders/12345678903
```
