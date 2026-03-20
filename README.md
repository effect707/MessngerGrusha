# Grusha Messenger

Мессенджер с поддержкой личных и групповых чатов, каналов, обмена файлами и уведомлений в реальном времени.

## Стек технологий

**Backend:** Go, gRPC, gRPC-Gateway (REST API), WebSocket

**Frontend:** React, TypeScript, Vite, Zustand

**Хранилища:** PostgreSQL, Redis, MinIO (S3-совместимое файловое хранилище)

**Инфраструктура:** Docker Compose, GitHub Actions CI, goose (миграции), sqlc, buf (protobuf)

## Архитектура

Проект построен по принципам Clean Architecture:

```
cmd/grusha/          — точка входа
internal/
  domain/            — бизнес-сущности
  usecase/           — бизнес-логика
  repository/        — слой данных (postgres, redis, minio)
  transport/         — транспортный слой (grpc, http, websocket)
  config/            — конфигурация через переменные окружения
  pkg/               — внутренние пакеты (jwt, hasher, pagination)
api/proto/           — protobuf-схемы
migrations/          — SQL-миграции (goose)
web/                 — React-клиент
```

## Возможности

- Регистрация и авторизация (JWT access + refresh токены)
- Личные и групповые чаты
- Каналы с ролевой моделью (owner, admin, member)
- Сообщения с пагинацией по курсору
- Реакции на сообщения
- Загрузка и скачивание файлов (MinIO)
- Уведомления в реальном времени через WebSocket
- Индикатор набора текста и онлайн-статус (Redis)

## Быстрый старт

### Запуск инфраструктуры

```bash
docker compose -f deployments/docker-compose.yml up -d
```

Поднимутся PostgreSQL, Redis и MinIO.

### Миграции

```bash
export DATABASE_URL="postgres://grusha:grusha_secret@localhost:5432/grusha?sslmode=disable"
make migrate-up
```

### Запуск сервера

```bash
make run
```

gRPC-сервер стартует на порту `50051`, HTTP/WS — на `8080`.

### Запуск фронтенда

```bash
cd web
npm install
npm run dev
```

## API

REST API доступен через gRPC-Gateway на `/api/`:

| Метод  | Эндпоинт                        | Описание                  |
|--------|----------------------------------|---------------------------|
| POST   | `/api/auth/register`             | Регистрация               |
| POST   | `/api/auth/login`                | Вход                      |
| POST   | `/api/auth/refresh`              | Обновление токенов        |
| GET    | `/api/chats`                     | Список чатов              |
| POST   | `/api/chats/direct`              | Создать личный чат        |
| POST   | `/api/chats/group`               | Создать групповой чат     |
| GET    | `/api/messages/{chat_id}`        | История сообщений         |
| POST   | `/api/messages`                  | Отправить сообщение       |
| GET    | `/api/channels`                  | Список каналов            |
| POST   | `/api/files/upload`              | Загрузить файл            |
| GET    | `/api/files/download?id=...`     | Скачать файл              |
| WS     | `/ws`                            | WebSocket-соединение      |

## Переменные окружения

| Переменная       | По умолчанию              | Описание                  |
|------------------|---------------------------|---------------------------|
| `GRPC_PORT`      | `50051`                   | Порт gRPC-сервера         |
| `HTTP_PORT`      | `8080`                    | Порт HTTP/WS-сервера      |
| `PG_HOST`        | `localhost`               | Хост PostgreSQL           |
| `PG_PORT`        | `5432`                    | Порт PostgreSQL           |
| `PG_USER`        | `grusha`                  | Пользователь БД           |
| `PG_PASSWORD`    | `grusha_secret`           | Пароль БД                 |
| `PG_DBNAME`      | `grusha`                  | Имя базы данных           |
| `REDIS_ADDR`     | `localhost:6379`          | Адрес Redis               |
| `MINIO_ENDPOINT` | `localhost:9000`          | Адрес MinIO               |
| `JWT_SECRET`     | `super-secret-key-change-me` | Секрет для JWT         |
| `JWT_ACCESS_TTL` | `15m`                     | Время жизни access-токена |
| `JWT_REFRESH_TTL`| `720h`                    | Время жизни refresh-токена|

## Makefile

```bash
make build            # Собрать бинарник
make run              # Запустить сервер
make test             # Запустить тесты
make lint             # Линтинг (golangci-lint)
make generate         # Сгенерировать protobuf + sqlc
make migrate-up       # Применить миграции
make migrate-down     # Откатить миграции
make docker-up        # Поднять инфраструктуру
make docker-down      # Остановить инфраструктуру
make web-dev          # Запустить фронтенд (dev)
make web-build        # Собрать фронтенд (prod)
```

## Тесты

```bash
make test
```

Unit-тесты покрывают слой usecase для всех модулей: auth, chat, message, channel, reaction, notification, file.
