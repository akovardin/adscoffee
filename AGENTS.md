# AGENTS.md

## Project overview

Go 1.27 monorepo (`go.ads.coffee/platform`) — an open-source advertising platform ("Caffeine"). Three binaries wired with `uber-go/fx` DI and `urfave/cli/v3`:

| Binary | Package | Default port | Purpose |
|---|---|---|---|
| admin | `admin/cmd/admin` | 8082 (container) / 8072 (published) | Admin panel (QOR5 UI) |
| server | `server/cmd/server` | 8081 (container) / 8071 (published) | Ad selection server |
| analytics | `analytics/cmd/analytics` | 8080 (dagu) / 8088 (published) | ClickHouse aggregation worker (dagu-scheduled) |

`dist/` is a self-contained production package (Docker Compose) that runs the whole platform: app services, Postgres, Redis Cluster, Kafka (3-node KRaft + `kafka-init`), ClickHouse, Kafka UI, Grafana, ch-ui, migrations, ClickHouse init SQL, dagu DAGs.

## Essential commands

```bash
# Local dev
task docker-up       # Start dev infrastructure (postgres, redis, kafka, clickhouse, minio, etc.)
task docker-down     # Stop dev infrastructure
task test            # go test -v ./...
task lint            # golangci-lint run ./...
go run ./admin/cmd/admin/main.go serve -config=./admin/configs/config.yaml     # Start admin locally
go run ./server/cmd/server/main.go serve -config=./server/configs/config.yaml  # Start server locally

# dist / production
task dist-update         # Sync dist assets from repo (redis-cluster, ClickHouse SQL, dagu dags, Grafana dashboards)
task dist-deploy         # Upload dist to the deploy host + up
task server-deploy       # Build + push + update only `server` (no full stack restart)
task admin-deploy        # Same for `admin` (also runs migrations)
task analytics-deploy    # Same for `analytics`
task nginx-deploy        # Upload nginx site configs + reload
task certbot-setup       # Install certbot, reload hook, enable auto-renewal
task certbot-renew-check # certbot renew --dry-run
```

## Dist & production deployment

- **Configs** live in `dist/{server,admin,analytics}/configs/config.yaml` and are the source of truth for production. They use `${VAR}` expansion (every service loads config via `go.uber.org/config` with `config.Expand(os.LookupEnv)`), so all hosts/credentials come from the environment.
- **`.env`** (`dist/.env`, gitignored; template `dist/.env.example`) holds secrets and service hosts (`POSTGRES_HOST`, `REDIS_ADDRS`, `KAFKA_SEEDS`, `CLICKHOUSE_HOST`, `S3_*`, `CookieSecure`, `GRAFANA_URL`, image tags, …). Docker Compose interpolates it and passes it to containers via `env_file`.
- **`.env.dev`** (`dist/.env.dev`, gitignored) holds `localhost` hosts for running binaries on the host against the Compose infra; used by VS Code `-dist` configs (`.vscode/launch.json`).
- **`.env.deploy`** (repo root, gitignored; template `.env.deploy.example`) holds `VENATOR=root@<host>` used by Taskfile `scp`/`ssh`. The `dist/*_IMAGE_TAG` variables choose the registry tags.
- **Taskfile/Supfile deploy to a remote host** (`venator`). `Supfile.yml` defines the `venator` network via the SSH alias `venator` (its `HostName`/`User` live in `~/.ssh/config`, not in git); Taskfile invokes `sup -sshconfig ~/.ssh/config`. `task dist-upload` rsyncs `dist/` to `/opt/platform`, `task dist-prepare` creates/chowns `data/` subdirs.
- **Migrations**: the one-shot `migrate` service applies schema; `admin` waits for it, `server` waits for `admin`. `server`/`clickhouse` also wait for `kafka-init` (which creates topics `request`, `impression`, `click`, `response`, `conversion`, `win`).
- **nginx / TLS**: site configs in `deploy/nginx/<domain>.coffee` (gitignored via `deploy/*`) reverse-proxy to the published localhost ports; certbot (webroot) renews certs, deploy-hook reloads nginx.

## QOR5 — local fork (must not change replace directives)

The admin UI uses a **locally-vendored fork** of QOR5. The `go.mod` at root has:

```
replace (
    github.com/qor5/admin/v3 => ./qor5/admin
    github.com/qor5/web/v3   => ./qor5/web
    github.com/qor5/x/v3     => ./qor5/x
)
```

Do **not** remove or redirect these `replace` directives. Any changes to QOR5 behavior must go into `qor5/`.

## Admin

`admin/cmd/admin` commands: `serve` (run), `migrate` (`m`, apply schema), `user` (`u`, create/update an admin user). `user` supports `-account`/`-a` and `-password`/`-p` (upsert by account). Grafana dashboard embedding uses `grafana.url` from config (`GRAFANA_URL` env).

## Server pipeline architecture

The ad server is a **plugin pipeline** (not middleware). Config in `server/configs/config.yaml` defines named pipelines as chains:

```
Route (e.g. /dsp/{placement}) -> Input -> Stage* -> Output
```

Each stage can declare `WithTargetings` to filter candidates. The central object is `*tools.State` (`server/internal/tools/state.go`) which carries the request, candidates, winners, placement, etc.

Plugin interfaces live in `server/internal/domain/plugins/`. Implementations in `server/plugins/`.

Important stages — `banners` (loads from DB/Redis cache), `limits` (budget/capping checks), `targeting` (filters candidates), `rotation` (weighted random pick), `mediation` (native-only, waterfall auction for third-party networks).

## Analytics pipeline

```
Server events -> Kafka topics -> ClickHouse raw tables (Kafka Engine) -> dagu aggregations (cron @ *:01) -> ClickHouse hourly tables -> Grafana / admin
```

Analytics aggregation is run by dagu (workflow scheduler). DAG YAMLs in `analytics/dagu/dags/` (copied into `dist/`). The binary command is `analytics aggregations -table=<requests|impressions|clicks|responses>`. ClickHouse also keeps raw tables for `conversion` and `win` topics (`analytics.*_hour` aggregates).

## Cached repositories

Server repos (`server/internal/repos/`) hold cached data (banners, placements, units) loaded from Postgres and stored in Redis. Three background goroutines warm these caches on startup.

## Configuration

Each service loads YAML config via `go.uber.org/config` (with `${VAR}` env expansion):

- Local dev: `admin/configs/config.yaml`, `server/configs/config.yaml` (committed); `dev.yaml`/`prod.yaml` are gitignored overrides.
- Production: `dist/{server,admin,analytics}/configs/config.yaml` (committed in `dist/`), values injected from `dist/.env`.

## Infrastructure (local dev)

`task docker-up` starts via `dev-compose.yaml`:
- Postgres (admin_dev, port 6432, user admin/123)
- Redis cluster (ports 7001-7006)
- Kafka (3-node KRaft, ports 19092/29092/39092)
- ClickHouse (ports 8123/9000, user admin/password)
- MinIO S3 (port 9222 API, 9223 console)
- Kafka UI (port 8090), ClickHouse UI (port 3488)

Most Go tests will fail unless these services are running.

## CI

Single workflow at `.github/workflows/go.yml`: runs on push/PR to `master`. Executes `golangci-lint` -> `go build -v ./...` -> `go test -v ./...` -> coverage report.

## Vendored dependencies

Root `vendor/` contains Go dependencies. The server also has a separate `server/vendor/`. Builds may use either module cache or vendored deps — check `GOFLAGS=-mod=vendor` if not set.

## Conventions

- Go code uses `golangci-lint` (no project-specific `.golangci.yml` at root — all config in vendored deps only)
- `go test -v ./...` (verbose) in all commands
- Docker images pushed to `registry.gitflic.ru/project/kovardin/adscoffee/platform-<service>`
- Deploy: `task <service>-deploy` (build -> publish -> update the single service on the remote host without stopping the stack) or `task dist-deploy` (upload whole `dist/` + up)
- Gitignored: `dist/.env`, `dist/.env.dev`, `.env.deploy`, `deploy/*`, `data/*`
