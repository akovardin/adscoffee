---
weight: 10
bookFlatSection: false
title: "Аналитика и мониторинг"
---

# 10 Аналитика и мониторинг
<p>&nbsp;</p>

Платформа собирает события показов, кликов и конверсий и показывает статистику в реальном времени.

## Как устроена аналитика

```
События сервера ──► Kafka ──► ClickHouse (сырые таблицы)
                                      │
                            dagu-агрегации (раз в час)
                                      │
                                      ▼
                       ClickHouse (*_hour) ──► Grafana и админка
```

Рекламный сервер публикует события в Kafka. ClickHouse читает их через Kafka Engine и складывает в сырые таблицы `requests`, `impressions`, `clicks`, `responses`, а также `conversion` и `win`. Сервис **analytics** (dagu) раз в час (`1 * * * *`) пересчитывает часовые витрины `analytics.requests_hour`, `impressions_hour`, `clicks_hour` и `responses_hour`. Для `conversion` и `win` сырые таблицы есть, но часовых DAG-ов нет.

Каждое событие несёт поля `network` и `unit_id` (победивший юнит рекламной сети) и денежное поле `revenue` — доход за показ. Именно `sum(revenue)` используется как денежная метрика.

## Дашборды

Аналитика в админке — это встроенный дашборд **Grafana** (раздел «Аналитика»): админка открывает iframe с настроенным дашбордом. URL Grafana задаётся переменной `GRAFANA_URL`.

Полноценные дашборды живут в **Grafana** (порт `3001`). В репозитории есть готовые дашборды:

- `analytics` — показы, клики, конверсии и доход;
- `platform-server` — метрики рекламного сервера;
- `platform-logs` — логи.

## ClickHouse

Для ручных запросов и просмотра таблиц есть **ch-ui** (порт `3488`), который подключается к ClickHouse внутри compose-сети. Креды ClickHouse задаются переменными `CLICKHOUSE_*` в `dist/.env`.

## Метрики и логи

- **Prometheus** (`9090`) собирает метрики сервисов.
- **Loki** (`3100`) хранит логи; **Promtail** собирает их с docker-контейнеров.
- Дашборды по метрикам и логам также доступны в Grafana.

## Обслуживание

Агрегации можно запускать вручную в интерфейсе dagu (порт `8088`, логин — `DAGU_AUTH_BASIC_USERNAME` / `DAGU_AUTH_BASIC_PASSWORD`). DAG-и лежат в `dist/analytics/dagu/dags/*.yaml`, конфиг аналитики — `dist/analytics/configs/config.yaml`.
