---
weight: 3
bookFlatSection: false
title: "Конфигурация"
---

# 03 Конфигурация
<p>&nbsp;</p>

Все параметры окружения задаются в файле `dist/.env`. Скопируйте шаблон и заполните значения:

```sh
cp dist/.env.example dist/.env
```

Конфиги `dist/{server,admin,analytics}/configs/config.yaml` являются источником истины — они подставляют значения из окружения через `${VAR}`. Docker Compose также прокидывает `.env` в контейнеры через `env_file`.

## Подключение к сервисам

| Переменная | Описание |
|---|---|
| `POSTGRES_HOST` / `POSTGRES_PORT` / `POSTGRES_USER` / `POSTGRES_PASSWORD` / `POSTGRES_DB` | Подключение к БД. Для внешней БД укажите её адрес и не поднимайте локальный `postgres` |
| `REDIS_ADDRS` | Адреса Redis Cluster, JSON-массивом (напр. `["redis:7001"]`) |
| `KAFKA_SEEDS` | Сиды Kafka, JSON-массивом (напр. `["kafka1:9092","kafka2:9092","kafka3:9092"]`) |
| `CLICKHOUSE_HOST` / `CLICKHOUSE_PORT` / `CLICKHOUSE_USER` / `CLICKHOUSE_PASSWORD` / `CLICKHOUSE_DB` | Подключение к ClickHouse |
| `S3_BUCKET` / `S3_ENDPOINT` / `S3_ACCESS_ID` / `S3_ACCESS_KEY` / `S3_REGION` / `S3_FORCE_PATH_STYLE` | Объектное хранилище (Yandex Object Storage, MinIO и т. п.). `S3_FORCE_PATH_STYLE=true` для MinIO |

## Интерфейсы и логины

| Переменная | Описание |
|---|---|
| `DAGU_AUTH_BASIC_USERNAME` / `DAGU_AUTH_BASIC_PASSWORD` | Логин и пароль dagu |
| `KAFKA_UI_USER` / `KAFKA_UI_PASSWORD` | Логин и пароль Kafka UI |
| `GRAFANA_ADMIN_USER` / `GRAFANA_ADMIN_PASSWORD` | Администратор Grafana |
| `GRAFANA_URL` | Внешний адрес Grafana для iframe в админке (должен открываться из браузера) |
| `GRAFANA_ROOT_URL` | Корневой URL Grafana (для корректных редиректов за reverse-proxy) |
| `GRAFANA_ANONYMOUS_ENABLED` | `false` — Grafana требует авторизацию; `true` — анонимный доступ (Viewer) |
| `CookieSecure` | `false` — админка открыта по plain HTTP (`http://IP:8072`); `true` — за HTTPS (иначе браузер не сохранит cookie и вход не сработает) |

## Телеметрия и образы

| Переменная | Описание |
|---|---|
| `TELEMETRY_ENABLED` | Включить трассировку сервера через OpenTelemetry (OTLP/gRPC) |
| `TELEMETRY_ENDPOINT` | Адрес OTLP-коллектора (напр. `jaeger:4317`) |
| `SERVER_IMAGE_TAG` / `ADMIN_IMAGE_TAG` / `ANALYTICS_IMAGE_TAG` | Теги образов сервисов из реестра |

## Важные замечания

- По умолчанию в `.env.example` заданы адреса для внутренней compose-сети (`postgres`, `redis`, `kafka1` и т. д.). Для запуска бинарников на хосте против поднятой инфраструктуры создайте `dist/.env.dev` с `localhost`-хостами (шаблон описан в комментариях `.env.example`).
- Для деплоя через Taskfile используется отдельный файл `.env.deploy` в корне репозитория с переменной `VENATOR=root@<host>` (шаблон — `.env.deploy.example`).
- Секреты из `.env` не нужно коммитить — файл исключён из git.
- Для HTTPS за reverse-proxy выставьте `CookieSecure=true`, иначе вход в админку не сохранится.
