# Setup guide

## Prerequisites

| Tool | Required | Notes |
| --- | --- | --- |
| Docker (+ Compose v2) | yes | Everything can run in containers |
| Make | yes | Every workflow is a Make target |
| Go 1.26+ | optional | Only for running the API directly on the host |
| Node 22+ | optional | Only for running Angular directly on the host |

When Node is missing, `frontend/Makefile` transparently runs every npm command
inside a `node:22-alpine` container, so `make test` and `make build` work either
way.

## Run the whole stack

```bash
make up
```

That command:

1. creates `infra/.env.local` from `infra/.env.example` on first run;
2. validates that every variable in the template is present and non-empty;
3. builds the backend and frontend images;
4. starts PostgreSQL, waits for it to be healthy, then the API, then nginx;
5. runs `infra/scripts/smoke.sh`, which checks `/healthz`, `/readyz`, the seeded
   catalog, and that the frontend served its runtime `config.json`.

Stop it with `make down`. The database volume survives; `make down` followed by
`docker volume rm midepensa_postgres_data` wipes it.

## Development loop

```bash
make dev
```

Starts PostgreSQL in Docker, then runs the API with hot reload and the Angular
dev server in parallel.

| Piece | URL | Reloads on |
| --- | --- | --- |
| API | <http://localhost:8080> | `.go` and `.sql` changes (via `air`) |
| Frontend | <http://localhost:4200> | any file under `src/` |

Install `air` and `golangci-lint` once with `make -C backend deps`. Without
`air`, `make dev` falls back to `go run`, which needs a manual restart.

### Backend only

```bash
cd backend
cp .env.example .env      # first time only
make services-start       # PostgreSQL on :5433
make dev                  # API on :8080
```

### Frontend only

```bash
cd frontend
make dev                  # Angular dev server on :4200
```

The dev server reads `public/config.json`, which points at
`http://localhost:8080/v1`. Change that file to develop against another API.

## Configuration

Configuration is environment variables everywhere. There is no config file
checked into the image.

### Backend (`backend/.env.example`)

| Variable | Meaning |
| --- | --- |
| `BACKEND_HOST`, `BACKEND_PORT` | Listen address |
| `AUTO_MIGRATE` | Apply pending migrations during startup |
| `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME` | PostgreSQL connection |
| `CORS_ORIGINS` | Comma-separated origin allowlist; `*` is local only |
| `LOG_LEVEL`, `LOG_FORMAT` | `debug\|info\|warn\|error`, `text\|json` |

`config.Load()` panics when a required variable is missing. That is deliberate:
a misconfigured deployment must fail at boot, not on the first request.

### Frontend

The frontend has **no build-time environment**. `entrypoint.sh` writes
`config.json` into the nginx root when the container starts, and
`ConfigService` fetches it before Angular bootstraps. One image, any backend.

| Variable | Meaning |
| --- | --- |
| `BACKEND_URL` | Absolute API base URL **as seen by the browser** |

### Stack (`infra/.env.example`)

`infra/.env.example` is the single source of truth: `validate_env.sh` fails a
deploy if `.env.<environment>` is missing any of its keys.

## Testing

```bash
make test                            # everything

make -C backend test-unit            # no dependencies
make -C backend test-integration     # needs `make -C backend services-start`
make -C frontend test                # Jest
```

Integration tests run against a real PostgreSQL, apply the migrations, and clean
up the rows they create.

## Troubleshooting

**`database files are incompatible with server`** — a Docker volume from another
project is being reused. All compose files here pin an explicit project name
(`name: midepensa-backend`), so this should not happen; if it does, remove the
stale volume with `docker volume ls` and `docker volume rm`.

**Backend container never becomes healthy** — check `docker logs
midepensa_local_backend`. The healthcheck must use `wget -O /dev/null` (a GET);
`wget --spider` sends a HEAD request, which Gin does not route.

**`config.json: BACKEND_URL must be an absolute http(s) URL`** — the frontend
container was started without `BACKEND_URL`, or `public/config.json` is wrong in
dev.

**Port already in use** — change `BACKEND_PORT` / `FRONTEND_PORT` in
`infra/.env.local`.
