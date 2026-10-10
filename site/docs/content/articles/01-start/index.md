---
weight: 1
bookFlatSection: false
title: "Быстрый старт"
---

# 01 Быстрый старт
<p>&nbsp;</p>

Платформа распространяется как исходный код в репозитории [github.com/akovardin/adscoffee](https://github.com/akovardin/adscoffee). Для запуска всей платформы разом используется каталог `dist/` — самодостаточный пакет на базе Docker Compose.

## Требования

- Docker Engine 24+ и плагин `docker compose` v2;
- сервер с доступом в интернет (виртуальная машина или физический хост).

## Установка

1. Клонируйте репозиторий и перейдите в каталог `dist/`:

   ```sh
   git clone https://github.com/akovardin/adscoffee.git
   cd adscoffee/dist
   ```

2. Создайте файл окружения и заполните секреты:

   ```sh
   cp .env.example .env
   ```

   Полный список переменных и их назначение описаны в разделе [Конфигурация](/docs/articles/03-config/).

3. Поднимите стек:

   ```sh
   docker compose up -d
   ```

   При первом запуске соберётся образ Redis Cluster и применятся миграции БД одноразовым сервисом `migrate`. Порядок старта: `postgres`/`redis`/`kafka` → `kafka-init` (создаёт топики) → `migrate` → `admin` → `server`; `analytics` поднимается после ClickHouse.

4. Создайте пользователя админки:

   ```sh
   docker compose run --rm admin user -c=/configs/config.yaml
   ```

   По умолчанию создаётся пользователь `admin` / `password` — обязательно смените пароль. Подробнее в разделе [Пользователи](/docs/articles/04-users/).

Применить миграции повторно (после обновления) можно так:

```sh
docker compose run --rm migrate
```

После этого админ-панель будет доступна по адресу `http://<ваш-хост>:8072`, а рекламный сервер — на порту `8071`.

## Состав пакета

В `dist/` входят:

- сервисы приложений: `server`, `admin`, `analytics` (на базе dagu);
- инфраструктура: PostgreSQL, Redis Cluster, Kafka, ClickHouse, Kafka UI, Grafana, ch-ui, Prometheus, Loki;
- миграции БД, init-скрипты ClickHouse, дашборды и датасорсы Grafana, DAG-и dagu;
- конфиги сервисов `dist/{server,admin,analytics}/configs/config.yaml`.

Подробнее о компонентах — в разделе [Архитектура](/docs/articles/02-architecture/).

## Доступные порты

| Сервис | Порт(ы) | Назначение |
|---|---|---|
| admin | 8072 | Админ-панель (QOR5) |
| server | 8071 | Рекламный сервер (`/dsp`, `/native`, `/banner`, `/inapp`, `/tracker`, `/postback`) |
| analytics (dagu) | 8088 | Планировщик агрегаций |
| grafana | 3001 | Дашборды |
| ch-ui | 3488 | Веб-клиент ClickHouse |
| kafka-ui | 8090 | Kafka UI |
| clickhouse | 8123 / 9000 | HTTP / native |
| postgres | 6432 | PostgreSQL |
| redis | 7001–7006 | Redis Cluster |
| kafka | 19092, 29092, 39092 | Внешние listener-порты брокеров |
| prometheus | 9090 | Метрики |
| loki | 3100 | Логи |

## Обновление и остановка

```sh
docker compose down        # остановить стек (данные в dist/data сохраняются)
docker compose up -d       # запустить снова
docker compose pull        # обновить образы сервисов
docker compose run --rm migrate   # применить миграции БД
```

## Деплой через Taskfile

В репозитории есть `Taskfile.yml` для сборки образов, публикации в реестр и обновления сервисов на удалённом хосте. Хост задаётся SSH-алиасом `VENATOR` в `.env.deploy` (шаблон — `.env.deploy.example`), сборка использует `Supfile.yml`.

```sh
task dist-update     # синхронизировать dist-ассеты из репозитория (dagu DAG-и, SQL, дашборды, redis-cluster)
task dist-deploy     # загрузить dist/ на хост и поднять стек
task server-deploy   # собрать/опубликовать и пересоздать только сервис server
task admin-deploy    # то же для admin (также выполняет миграции)
task analytics-deploy

task nginx-deploy    # залить конфиги nginx и перезагрузить
task certbot-setup   # установить certbot и включить автопродление
```

`task <service>-deploy` не синхронизирует `dist/` и не трогает конфиги сервисов: конфиги монтируются с хоста (`./{server,admin,analytics}/configs`), поэтому изменения в `config.yaml` вступят в силу только после `task dist-upload` и пересоздания сервиса.

Следующие шаги: настройте [конфигурацию](/docs/articles/03-config/), создайте [пользователя](/docs/articles/04-users/) и заведите первый [сайт](/docs/articles/05-sites/).
