---
weight: 2
bookFlatSection: false
title: "Архитектура"
---

# 02 Архитектура
<p>&nbsp;</p>

Платформа — это Go-монорепозиторий из трёх приложений и набора инфраструктурных сервисов. Стек поднимается одним `docker compose up -d` из каталога `dist/`.

## Приложения

### Админка (`admin`)

Веб-панель управления на QOR5. Через неё создаются и настраиваются рекламодатели, кампании, группы, креативы, рекламные сети, сайты, плейсменты, юниты и пользователи, а также отображается встроенный дашборд Grafana. Загруженные изображения хранятся в S3-совместимом объектном хранилище.

- Порт: `8072` (`8082` внутри контейнера).
- Конфиг: `dist/admin/configs/config.yaml`.
- Команды образа: `serve` (запуск), `migrate` (миграции БД), `user` (создать пользователя или сменить пароль).
- Зависит от PostgreSQL, ClickHouse и S3.

### Рекламный сервер (`server`)

Сервер показа рекламы. Работает как конвейер плагинов `Route → Input → Stages → Output`: загрузка кандидатов, проверка лимитов и бюджетов, таргетинг, взвешенный выбор (rotation) и медиация. Сервер не хранит состояние в процессе, поэтому его можно масштабировать горизонтально.

- Порт: `8071` (`8081` внутри).
- Конфиг: `dist/server/configs/config.yaml` — здесь описаны пайплайны маршрутов `/dsp`, `/native`, `/banner`, `/inapp`, `/tracker`, `/postback`.
- Использует PostgreSQL (кэшируемые репозитории баннеров, плейсментов и юнитов), Redis Cluster (кэш) и публикует события в Kafka.

### Аналитика (`analytics`)

Планировщик и воркер агрегаций на базе **dagu**. По расписанию раз в час (`1 * * * *`) запускает агрегацию данных в ClickHouse по таблицам `requests`, `impressions`, `clicks`, `responses` и складывает результат в часовые витрины `analytics.*_hour` для дашбордов. События `conversion` и `win` хранятся в сырых таблицах, часовых DAG-ов для них нет.

- Порт: `8088` (`8080` внутри).
- Образ: `platform-analytics`; команда агрегации — `analytics aggregations -table=<requests|impressions|clicks|responses> -c=/configs/config.yaml`.
- Конфиг: `dist/analytics/configs/config.yaml`; DAG-и — `dist/analytics/dagu/dags/*.yaml`, данные и логи — `dist/data/dagu/{data,logs}`.
- Читает события из Kafka через Kafka Engine ClickHouse, пишет агрегаты в `analytics.*_hour` (включая `revenue`, `network`, `unit_id`, `bundle`).

## Поток событий

```
Рекламный сервер ──► Kafka (request, impression, click, response, conversion, win)
                              │
                              ▼
                    ClickHouse (сырые таблицы)
                              │
                    dagu-агрегации (раз в час)
                              │
                              ▼
                 ClickHouse (*_hour) ──► Grafana / админка
```

## Инфраструктура

- **Kafka** — шина событий платформы: кластер из трёх брокеров в режиме KRaft (`kafka1`, `kafka2`, `kafka3`, `kafka-init`, `kafka-ui`). Топики создаются сервисом `kafka-init`. Kafka UI — на порту `8090`.
- **PostgreSQL** — основное хранилище данных админки и сервера. Порт `6432`.
- **ClickHouse** — хранилище событий и агрегатов для аналитики. Порты `8123` (HTTP) и `9000` (native).
- **Redis Cluster** — кэш сервера показа. Порты `7001–7006`.
- **S3** — объектное хранилище для загружаемых изображений (Yandex Object Storage, MinIO, AWS S3 и другие).
- **Grafana** — дашборды. Порт `3001`.
- **ch-ui** — веб-клиент ClickHouse для ручных запросов. Порт `3488`.
- **Prometheus и Loki** — метрики и логи (порты `9090` и `3100`), собираются через Promtail.

## SDK

Для интеграции рекламы в приложения и на сайты есть SDK:

- **Web / JS** — `server/static/js/coffee.js` (`window.adsCoffee.render(...)`), поддерживает собственные креативы и блоки Yandex.
- **Android** — `sdk/android/` (модуль `coffeesdk`): баннеры, интерстишелы, форматы coffee и Yandex.

Денежное поле событий называется `revenue` (не `price`) на всём пути: сервер → Kafka → ClickHouse (`analytics.*.revenue`, `*_hour.revenue`). Ставка `Banner.Price`/`Unit.Price` (CPM, выбор победителя) остаётся с именем `price`.

Подробнее — в разделе [Интеграция рекламы](/docs/articles/09-integration/).
