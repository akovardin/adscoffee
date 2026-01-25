# Вектор

Рекламный сервер, который можно бесконечно расширять плагинами

Статья про различные системы расширений

- https://www.youtube.com/watch?v=qerWv9JTlo8&list=PLJTW0ZQ22rrFwoMuJdyIhtryfyJYghrvd&index=3
- https://github.com/GopherConRu/talks/blob/main/2020/Designing%20Pluggable%20Idiomatic%20Go%20Applications%20-%20Mark%20Bates.pdf

Данные в стейте прокидываются через контекст. Это поможет прокинуть параметры из входящего плагина в исходщие плагины

## Caddy

https://caddyserver.com/docs/extending-caddy

## hashicorp

https://github.com/hashicorp/go-plugin

## Telegraf

Типы плагинов
- Input
- Output
- Aggregator (filter?)
- Processor

- https://github.com/influxdata/telegraf
- https://docs.influxdata.com/telegraf/v1/plugins/


## Jaeger

Типы плагинов
- Storage
- Query
- Sampling
- Exporters
- Processors

- https://www.jaegertracing.io/docs/2.11/deployment/configuration/#extensions


## Docker

- https://docs.docker.com/extensions/extensions-sdk/architecture/

Расширения для докера - это специальные Docker контейнер. В корне файловой системы образа находится файл metadata.json, который описывает содержимое расширения. Это фундаментальный элемент расширения Docker.

Расширение может содержать пользовательский интерфейс и серверную часть, которые работают либо на хосте, либо на виртуальной машине Desktop.

Сразу можно посмотреть на примеры расширений: https://github.com/docker/extensions-sdk/tree/main/samples

Docker предоставляет набор команд для создания расширений, например есть специальная команда, которая генерирует шаблон:

```
docker extension init <my-extension>
```

Эта команда сгенерирует шаблон, в котором будет:
- Сервер на Go в папке backend, который прослушивает сокет. У сервера есть одна конечная точка /hello, которая возвращает данные в формате JSON.
- Фронтенд на React в папке frontend, который может вызывать сервер и выводить ответ.

Кроме этого, есть еще несколько команд, которые помогут собрать расширение:

```
docker build -t <name-of-your-extension> .
```

И установить его

```
docker extension install <name-of-your-extension>
```

TODO: что доступно в расширении

Вы можете публиковаться в Extensions Marketplace

## Traefik

- https://plugins.traefik.io/create
- https://doc.traefik.io/traefik/extend/extend-traefik/

Как насчет https://github.com/traefik/yaegi ?

## Zabix

- https://habr.com/ru/companies/jetinfosystems/articles/963662/

## Grafa

- https://grafana.com/developers/plugin-tools/key-concepts/anatomy-of-a-plugin

## Hugo

- Modules https://gohugo.io/hugo-modules/introduction/


## Плагины для рекламной системы

### Вопросы

Как добавлять новый таргетинг?
Как пробрасывать новый тип баннера через всю структуру?

### Работа плагинов

Что имеет смыл убрать в плагины:

- Приемник запроса
- Этапы обработки запросаы
- Таргетинги. Каждый таргетинг в свой плагин
- Форматы ответа

Нужно объяснить концепцию пайплайнов. Обработка запроса, как правило, идет по заданному сценарию. Например, обрабатывается запрос от определенной SSP, дальше обрабатываются таргетинги, лимиты и нужно правильно сформировать ответ в зависимости от протокола. Эта логика образует своеобразны пайплайн

Пример конфигурации нескольких пайплайнов

```yaml
pipelines:

  - name: google
    input:
      name: rtb
    stages:
      - name: banners
      - name: limits
      - name: targeting
      - name: rotation
    targetings:
      - name: geo
      - name: ua    
    output:
      name: output
      config:
        formats: 
          - native
          - banner

  - name: amazon
    input:
      name: rtb
    stages:
      - name: banners
      - name: limits
      - name: targeting
      - name: rotation
      - name: process
    targetings:
      - name: geo
      - name: ua
    output:
      name: output
      config:
        formats: 
          - native
          - banner
```

Что важно - хочется собирать разные пайплайны. Те один плагин должен инстанционироваться больше одного раза. И тут нужно как-то пробросить зависимости, а это приводит нас к том, что нам нужны два экземпляра - билдер и сам плагин

Есть специальный класс - `pipeline.Manager`. Он получает на вход конфиг и набор менеджеров для каждого типа плагина. Менеджер в методе `NewManager` получает конфиг и формирует набор пайплайнов

Каждый плагин предоставляется через `fx`. Это значит, что в плагине есть доступ ко всем зависимостям. Но при формировании пайплайна к плагина вызывается метод `Build(cgf map[string]any) Plugin`, который заново создает экземпляр плагина, но уже с заданными настройками. Это нужно для того, чтобы в разных пайплайнах можно было настраивать плагин по разному

Формат ответа может зависеть от запроса. Например, в rtb в зависимости от запроса отдаются нужные форматы ответа. Для этого можно реализовать похожую логику как с таргетингами

Для большей гибкости должна автоматически заводиться ручка по названию пайплайна


## Выводы

Для системы плагинов не подойдет fx. Через fx можно подключать менеджеры плагинов и менеджер паплайнов чтобы сформировать набор пайплайнов в которых будут использоваться плагины 