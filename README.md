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

## Сборка
```bash
# Добавить go в PATH
source ~/.profile

# Запустить сервер
go build -o cmd/gophermart/gophermart cmd/gophermart/*go \
    && ./cmd/gophermart/gophermart
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
    --cookie "auth_token=TOKEN" \
    --data '12345678903' \
    --include \
    localhost:8081/api/user/orders 
```
