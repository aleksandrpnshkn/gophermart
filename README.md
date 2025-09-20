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

## Сборка
```bash
# Добавить go в PATH
source ~/.profile

# Запустить unit-тесты (count для отключения кэша, помогает отлавливать flaky-тесты)
go test -count=10 ./...

# Запустить сервер
go build -o cmd/gophermart/gophermart cmd/gophermart/*go \
    && ./cmd/gophermart/gophermart
```

## Тестовые запросы
```bash
curl --include localhost:8081/
curl --include localhost:8081/api/ping
```
