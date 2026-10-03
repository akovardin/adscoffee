# Caffeine

[![Go Build](https://github.com/akovardin/adscoffee/actions/workflows/go.yml/badge.svg)](https://github.com/akovardin/adscoffee/actions)
[![Go Coverage](https://github.com/akovardin/adscoffee/wiki/coverage.svg)](https://raw.githack.com/wiki/akovardin/adscoffee/coverage.html)

## Рекламная платформа для всех

Кофеин — это полнофункциональная рекламная платформа, созданная для того, чтобы дать разработчикам полный контроль над монетизацией своих проектов. Прозрачная, настраиваемая и этичная альтернатива крупным рекламным сетям.

Кофеин позволяет самостоятельно создавать и запускать свои рекламные кампании в своих приложениях. Вы можете работать с прямыми рекламодателями.

Философия проекта строится на нескольких фундаментальных принципах.

Прозрачность — это основа доверия. Весь исходный код платформы открыт для независимого аудита. Вы всегда будете точно знать, как работает алгоритм выбора рекламы, какие данные собираются и каким образом они обрабатываются. Никаких скрытых алгоритмов или «чёрных ящиков».

Контроль полностью в ваших руках. Вы самостоятельно управляете всеми аспектами монетизации: от дизайна и размещения рекламных блоков до тонкой настройки таргетинга и частоты показа. Платформа — это ваш инструмент, а не набор навязанных правил.

Этика лежит в основе подхода к пользователю. Система спроектирована с уважением к аудитории. Вы сможете легко устанавливать строгие правила к качеству и уместности рекламы, избегать навязчивых форматов и гарантированно соблюдать регуляторные нормы, такие как GDPR.

## Установка и конфигурация

Для запуска всей платформы разом используется каталог `dist/` — самодостаточный пакет на базе Docker Compose. В него входят:

- сервисы приложений: `server`, `admin`, `analytics` (dagu);
- инфраструктура: PostgreSQL, Redis Cluster, Kafka (3 узла, KRaft), ClickHouse, Kafka UI, Grafana, ch-ui;
- миграции БД (`migrate`), init-скрипты ClickHouse, дашборды/датасорс Grafana, DAG-и dagu;
- конфиги сервисов `dist/{server,admin,analytics}/configs/config.yaml`.

Требования: Docker Engine 24+ и плагин `docker compose` v2.

### Быстрый старт

1. Создать файл окружения и заполнить секреты:

   ```sh
   cp dist/.env.example dist/.env
   ```

2. Поднять стек:

   ```sh
   cd dist
   docker compose up -d
   ```

   При первом запуске соберётся образ Redis Cluster и применятся миграции. Порядок старта: `migrate` (создаёт таблицы) → `admin` → `server`; `analytics` поднимается после ClickHouse.

3. Создать пользователя админки (подробнее — в разделе «Пользователи админки»):

   ```sh
   docker compose run --rm admin user -c=/configs/config.yaml
   ```

### Пользователи админки

Управление пользователями выполняется командой `user` (алиас `u`) образа `platform-admin`. Команда создаёт пользователя с указанным логином, а если такой уже существует — обновляет ему пароль:

```sh
# в каталоге dist/
docker compose run --rm admin user -c=/configs/config.yaml

# свой логин и пароль
docker compose run --rm admin user -c=/configs/config.yaml -account=bob -password=S3cret
```

Для уже запущенного контейнера (например, на удалённом хосте) можно выполнить команду без создания нового контейнера:

```sh
docker compose exec admin /admin user -c=/configs/config.yaml -account=bob -password=S3cret
```

Флаги:

| Флаг | По умолчанию | Описание |
|---|---|---|
| `--account`, `-a` | `admin` | Логин пользователя |
| `--password`, `-p` | `password` | Пароль пользователя |
| `--config`, `-c` | — | Путь к конфигу внутри контейнера (`/configs/config.yaml`) |

Если конфиг не указан, используется `admin/configs/config.yaml`. Новых пользователей также можно добавлять прямо в админке — в разделе «Пользователи».

> По умолчанию создаётся `admin` / `password` — обязательно смените пароль. Флаги `-account`/`-password` доступны в образе, собранном из текущих исходников (`task admin-deploy` или `docker build -f admin/Dockerfile`).

### Доступные порты

| Сервис | Порт(ы) | Назначение |
|---|---|---|
| admin | 8072 | Админ-панель (QOR5) |
| server | 8071 | Рекламный сервер (`/dsp`, `/native`, `/banner`, `/tracker`, `/postback`) |
| analytics (dagu) | 8088 | Планировщик агрегаций |
| grafana | 3001 | Дашборды |
| ch-ui | 3488 | Веб-клиент ClickHouse |
| kafka-ui | 8090 | Kafka UI |
| clickhouse | 8123 / 9000 | HTTP / native |
| postgres | 6432 | PostgreSQL |
| redis | 7001–7006 | Redis Cluster |
| kafka | 19092, 29092, 39092 | Внешние listener-порты брокеров |

### Что нужно настроить

Все параметры задаются в `dist/.env`:

| Переменная | Описание |
|---|---|
| `POSTGRES_HOST` / `POSTGRES_PORT` / `POSTGRES_USER` / `POSTGRES_PASSWORD` / `POSTGRES_DB` | Подключение к БД. Для внешней БД укажите её адрес и не поднимайте локальный `postgres` |
| `REDIS_ADDRS` | Адреса Redis Cluster, JSON-массивом (напр. `["redis:7001"]`) |
| `KAFKA_SEEDS` | Сиды Kafka, JSON-массивом (напр. `["kafka1:9092","kafka2:9092","kafka3:9092"]`) |
| `CLICKHOUSE_HOST` / `CLICKHOUSE_PORT` / `CLICKHOUSE_USER` / `CLICKHOUSE_PASSWORD` / `CLICKHOUSE_DB` | Подключение к ClickHouse |
| `S3_BUCKET` / `S3_ENDPOINT` / `S3_ACCESS_ID` / `S3_ACCESS_KEY` / `S3_REGION` / `S3_FORCE_PATH_STYLE` | Объектное хранилище (Yandex Object Storage, MinIO, etc). `S3_FORCE_PATH_STYLE=true` для MinIO |
| `DAGU_AUTH_BASIC_USERNAME` / `DAGU_AUTH_BASIC_PASSWORD` | Логин dagu |
| `KAFKA_UI_USER` / `KAFKA_UI_PASSWORD` | Логин Kafka UI |
| `GRAFANA_ADMIN_USER` / `GRAFANA_ADMIN_PASSWORD` | Админ Grafana |
| `CookieSecure` | `false` — админка по plain HTTP (`http://IP:8072`); `true` — за HTTPS (иначе браузер не сохранит cookie и логин не сработает) |
| `SERVER_IMAGE_TAG` / `ADMIN_IMAGE_TAG` / `ANALYTICS_IMAGE_TAG` | Теги образов из `registry.gitflic.ru/project/kovardin/adscoffee` |

Конфиги `dist/{server,admin,analytics}/configs/config.yaml` — источник истины, они подставляют значения из окружения через `${VAR}`. Остальные данные (миграции ClickHouse, дашборды Grafana, DAG-и dagu, Redis Cluster) обновляются из репозитория командой:

```sh
task dist-update
```

### Отладка на хосте

Для запуска бинарников локально против инфраструктуры из Compose создайте `dist/.env.dev` с `localhost`-хостами (шаблон есть в `dist/.env.example`) — его используют конфигурации VS Code (`-dist` в `.vscode/launch.json`).

### Обновление и деплой

```sh
task dist-deploy        # залить dist на удалённый хост и запустить/обновить стек
task server-deploy      # собрать, запушить и обновить только server (без остановки стека)
task admin-deploy       # то же для admin (плюс прогон миграций)
task analytics-deploy   # то же для analytics
```

Остановка: `docker compose down` (данные в `dist/data` сохраняются).

## Компоненты системы

Стек поднимается одним `docker compose up -d` из `dist/` (см. `dist/docker-compose.yaml`). Основные сервисы:

### Админка (`admin`)
Веб-панель управления на QOR5. Через неё создаются и настраиваются рекламодатели, кампании, группы, баннеры, сети, площадки, юниты и пользователи, а также отображается встроенный дашборд Grafana.

- Порт: `8072` (`8082` внутри контейнера).
- Конфиг: `dist/admin/configs/config.yaml` (БД, ClickHouse, S3, URL Grafana).
- Команды образа: `serve` (запуск), `migrate` (миграции БД, сервис `migrate`), `user` (создать пользователя/сменить пароль).
- Зависит от PostgreSQL, ClickHouse, S3; логин через qor5 (`CookieSecure`).

### Рекламный сервер (`server`)
Сервер показа рекламы. Работает как конвейер плагинов `Route → Input → Stages → Output`: загрузка кандидатов, проверка лимитов/бюджетов, таргетинг, взвешенный выбор (rotation), медиация.

- Порт: `8071` (`8081` внутри).
- Конфиг: `dist/server/configs/config.yaml` (описание пайплайнов: `/dsp`, `/native`, `/banner`, `/inapp`, `/tracker`, `/postback`).
- Использует PostgreSQL (кэшируемые репозитории баннеров/плейсментов/юнитов), Redis Cluster (кэш) и публикует события в Kafka.

### Аналитика (`analytics`)
Планировщик и воркер агрегаций на базе **dagu**. По расписанию (`1 * * * *`) запускает агрегацию данных в ClickHouse по таблицам `requests`, `impressions`, `clicks`, `responses` и складывает результат в часовые витрины для дашбордов.

- Порт: `8088` (`8080` внутри).
- Конфиг: `dist/analytics/configs/config.yaml`; DAG-и — `dist/analytics/dagu/dags/*.yaml`.
- Читает события из Kafka через Kafka Engine ClickHouse, пишет агрегаты в `analytics.*_hour`.

### Kafka (`kafka1`, `kafka2`, `kafka3`, `kafka-init`, `kafka-ui`)
Шина событий платформы (кластер из 3 брокеров в режиме KRaft). Сервер публикует события, ClickHouse их читает и парсит.

- Топики (создаются сервисом `kafka-init`): `request`, `impression`, `click`, `response`, `conversion`, `win`.
- Внешние listener-порты: `19092`, `29092`, `39092`.
- Kafka UI: порт `8090` (`8080` внутри); логин `KAFKA_UI_USER`/`KAFKA_UI_PASSWORD` из `.env`.

### ch-ui (`ch-ui`)
Веб-клиент ClickHouse для ручных запросов и просмотра таблиц.

- Порт: `3488`.
- Подключается к `clickhouse:8123` внутри compose-сети.

### Дополнительные сервисы
Остальная инфраструктура стека: `postgres` (`6432`), `redis` (Redis Cluster, `7001–7006`), `clickhouse` (`8123`/`9000`), `grafana` (`3001`), а также SDK для показа рекламы на сайте и в Android.

## Масштабирование и промышленная эксплуатация

Система не привязана к инфраструктуре из `dist/`: все внешние зависимости задаются через окружение (`dist/.env`), поэтому любой компонент можно заменить на облачный или управляемый без изменения кода. Система готова к расширению в любой момент.

- **Kafka.** Вместо локального кластера из Docker достаточно указать адреса облачной (Managed) Kafka в `KAFKA_SEEDS`; конфиг поддерживает SASL (`sasl_mechanism`, `username`, `password`) и TLS.
- **PostgreSQL.** Выносится на отдельный сервер или в управляемую БД через `POSTGRES_HOST`, `POSTGRES_PORT`, `POSTGRES_USER`, `POSTGRES_PASSWORD`, `POSTGRES_DB` (локальный `postgres` при этом не поднимается).
- **ClickHouse.** Аналогично — внешний или облачный инстанс через `CLICKHOUSE_HOST`, `CLICKHOUSE_PORT`, `CLICKHOUSE_USER`, `CLICKHOUSE_PASSWORD`, `CLICKHOUSE_DB`.
- **Redis.** Внешний Redis Cluster через `REDIS_ADDRS`.
- **S3.** Любое S3-совместимое хранилище (Yandex Object Storage, MinIO, AWS S3) через `S3_*`.
- **Рекламный сервер.** `server` не хранит состояние в процессе (данные — в PostgreSQL/Redis, события — в Kafka), поэтому можно запускать сколько угодно инстансов за балансировщиком: на отдельных серверах или в Kubernetes, с горизонтальным масштабированием.
- **Админка и аналитика.** Масштабируются аналогично: `admin` — stateless (состояние в БД), `analytics` (dagu) — по нагрузке.

Миграции БД и init-скрипты ClickHouse идемпотентны, вся конфигурация — через переменные окружения. В итоге система готова к промышленному использованию и в любой момент может быть разнесена на несколько серверов или переведена в Kubernetes.

## Разработка

Для локальной разработки используется Docker Compose (инфраструктура) и [Taskfile](https://taskfile.dev/). Полный стек со всеми сервисами и мониторингом проще поднять из `dist/` — см. раздел «Установка и конфигурация».

1. Поднять инфраструктуру (PostgreSQL, Redis Cluster, Kafka, ClickHouse, MinIO, Grafana, Kafka UI):

   ```sh
   task docker-up
   ```

2. Запустить нужные сервисы платформы:

   ```sh
   # админка
   go run ./admin/cmd/admin/main.go serve -config=./admin/configs/config.yaml

   # рекламный сервер
   go run ./server/cmd/server/main.go serve -config=./server/configs/config.yaml
   ```

3. Применить миграции БД (создание таблиц):

   ```sh
   go run ./admin/cmd/admin/main.go migrate -config=./admin/configs/dev.yaml
   ```

4. Создать пользователя админки (логин и пароль задаются флагами):

   ```sh
   go run ./admin/cmd/admin/main.go user -account=admin -password=secret -config=./admin/configs/dev.yaml
   ```

Локальные конфиги — `admin/configs/config.yaml` и `server/configs/config.yaml`; персональные оверрайды `dev.yaml`/`prod.yaml` не коммитятся. В VS Code есть готовые конфигурации запуска (`.vscode/launch.json`).

Остановить инфраструктуру и дополнительные команды:

```sh
task docker-down   # остановить инфраструктуру
task test          # go test -v ./...
task lint          # golangci-lint run ./...
```

