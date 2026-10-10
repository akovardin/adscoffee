---
weight: 13
bookFlatSection: false
title: "Разработка"
---

# 13 Разработка
<p>&nbsp;</p>

Проект написан на Go и распространяется как открытый исходный код. Мы рады вкладу сообщества.

- Репозиторий: [github.com/akovardin/adscoffee](https://github.com/akovardin/adscoffee)
- Сборка и тесты: `.github/workflows/go.yml` (golangci-lint → build → test → coverage)

## Требования

- Go (см. `go.mod` в корне репозитория);
- Docker Engine и `docker compose` v2;
- [Taskfile](https://taskfile.dev/).

## Локальная разработка

1. Поднимите инфраструктуру (PostgreSQL, Redis Cluster, Kafka, ClickHouse, MinIO, Kafka UI):

   ```sh
   task docker-up
   ```

2. Запустите нужные сервисы:

   ```sh
   # админка
   go run ./admin/cmd/admin/main.go serve -config=./admin/configs/config.yaml

   # рекламный сервер
   go run ./server/cmd/server/main.go serve -config=./server/configs/config.yaml

   # аналитика
   go run ./analytics/cmd/analytics/main.go aggregations -table=impressions -config=./analytics/configs/config.yaml
   ```

3. Примените миграции БД:

   ```sh
   go run ./admin/cmd/admin/main.go migrate -config=./admin/configs/config.yaml
   ```

4. Создайте пользователя админки:

   ```sh
   go run ./admin/cmd/admin/main.go user -account=admin -password=secret -config=./admin/configs/config.yaml
   ```

Локальные конфиги — `admin/configs/config.yaml`, `server/configs/config.yaml`, `analytics/configs/config.yaml`. Персональные оверрайды `dev.yaml` / `prod.yaml` не коммитятся. В VS Code есть готовые конфигурации запуска (`.vscode/launch.json`).

## Полезные команды

```sh
task docker-up      # поднять инфраструктуру для разработки
task docker-down    # остановить инфраструктуру
task test           # go test -v ./...
task lint           # golangci-lint run ./...
task dist-update    # синхронизировать dist-ассеты (dagu DAG-и, SQL, дашборды, redis-cluster) из репозитория
```

## Структура репозитория

| Каталог | Назначение |
|---|---|
| `admin/` | Админ-панель (QOR5) |
| `server/` | Рекламный сервер (пайплайн плагинов) |
| `analytics/` | Планировщик агрегаций (dagu) и DAG-и |
| `qor5/` | Локальный форк QOR5 (`replace`-директивы в `go.mod`); изменения поведения QOR5 вносятся сюда |
| `pkg/` | Общие пакеты (redis, kafka, logger, telemetry и др.) |
| `sdk/` | SDK для интеграции (Android) |
| `geoip/` | Данные GeoIP для таргетинга по гео |
| `data/` | Миграции и init-скрипты ClickHouse |
| `docker/` | Docker-материалы (redis-cluster, grafana dashboards) |
| `dist/` | Самодостаточный пакет для запуска платформы |
| `deploy/` | Конфиги nginx для продакшена (не коммитятся) |

> `replace`-директивы QOR5 в `go.mod` переписывать нельзя — админка использует локальный форк.

## Как внести вклад

1. Сделайте форк репозитория и создайте ветку.
2. Внесите изменения, добавьте тесты.
3. Прогоните `task lint` и `task test`.
4. Откройте pull request с описанием изменений.

Перед изменениями ядра рекомендуется обсудить идею в issue — так проще договориться о подходе.
