---
weight: 12
bookFlatSection: false
title: "Эксплуатация"
---

# 12 Эксплуатация и масштабирование
<p>&nbsp;</p>

Платформа не привязана к инфраструктуре из `dist/`: все внешние зависимости задаются через переменные окружения (`.env`), поэтому любой компонент можно заменить на облачный или управляемый без изменения кода.

## Внешние зависимости

- **Kafka.** Вместо локального кластера из Docker достаточно указать адреса любых серверов Kafka в `KAFKA_SEEDS`; конфиг поддерживает SASL (`sasl_mechanism`, `username`, `password`) и TLS.
- **PostgreSQL.** Выносится так же, как Kafka, через `POSTGRES_HOST`, `POSTGRES_PORT`, `POSTGRES_USER`, `POSTGRES_PASSWORD`, `POSTGRES_DB`.
- **ClickHouse.** Внешний или облачный инстанс через `CLICKHOUSE_HOST`, `CLICKHOUSE_PORT`, `CLICKHOUSE_USER`, `CLICKHOUSE_PASSWORD`, `CLICKHOUSE_DB`.
- **Redis.** Внешний Redis Cluster через `REDIS_ADDRS`.
- **S3.** Любое S3-совместимое хранилище (Yandex Object Storage, MinIO, AWS S3) через переменные `S3_*`. Для MinIO установите `S3_FORCE_PATH_STYLE=true`.

## Горизонтальное масштабирование

- **Рекламный сервер** не хранит состояние в процессе: данные берутся из PostgreSQL и Redis, события летят в Kafka. Поэтому можно запускать сколько угодно инстансов сервера за балансировщиком — на отдельных серверах или в Kubernetes.
- **Админка** масштабируется по числу одновременных пользователей.
- **Аналитика** масштабируется по нагрузке и объёму данных.

Миграции БД и init-скрипты ClickHouse идемпотентны, вся конфигурация выполняется через переменные окружения. В итоге систему можно разнести на несколько серверов или перевести в Kubernetes.

## Обновление

```sh
cd dist
docker compose pull        # подтянуть новые образы сервисов
docker compose up -d       # применить обновление
docker compose run --rm migrate   # применить миграции БД
```

Для удалённого хоста в репозитории есть Taskfile: `task dist-deploy` (загрузить `dist/` и поднять стек) и `task {server,admin,analytics}-deploy` (собрать образ сервиса, опубликовать в реестр и пересоздать только его). Конфиги сервисов монтируются с хоста (`./{server,admin,analytics}/configs`), поэтому правки `config.yaml` применяются только после `task dist-update`/`dist-upload` и пересоздания сервиса.

Данные сохраняются в каталоге `dist/data` (PostgreSQL, ClickHouse, Grafana, Kafka, dagu и другие) и переживают перезапуск стека.

## Резервное копирование

Резервируйте как минимум:

- каталог `dist/data/postgres` — данные админки и рекламного сервера;
- каталог `dist/data/clickhouse` — события и агрегаты;
- S3-бакет с загруженными изображениями;
- файл `dist/.env` с секретами.

Для продакшена рекомендуется использовать управляемые PostgreSQL/ClickHouse/Redis/S3 и настраивать их штатные бэкапы.

## Безопасность

- Обязательно смените все пароли из `.env.example` (`POSTGRES_PASSWORD`, `CLICKHOUSE_PASSWORD`, `GRAFANA_ADMIN_PASSWORD`, `KAFKA_UI_PASSWORD`, `DAGU_AUTH_BASIC_PASSWORD`, `S3_*`).
- Не коммитьте `dist/.env`.
- Для HTTPS за reverse-proxy устанавливайте `CookieSecure=true`.
- Ограничьте доступ к инфраструктурным портам (ClickHouse, PostgreSQL, Redis, Kafka) из интернета.
