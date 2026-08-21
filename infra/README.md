# Infrastructure

One `docker-compose.yml` drives every environment. What changes between them is
the env file, never the compose file.

## Services

| Service | Image | Published on | Purpose |
| --- | --- | --- | --- |
| `postgres` | `${POSTGRES_IMAGE}` | internal only | Database |
| `backend` | `${BACKEND_IMAGE}` | `127.0.0.1:${BACKEND_PORT}` | Go API |
| `frontend` | `${FRONTEND_IMAGE}` | `127.0.0.1:${FRONTEND_PORT}` | nginx serving the Angular build |

They start in dependency order: nothing talks to a service that has not reported
healthy.

Ports bind to `127.0.0.1`, so outside local development a host reverse proxy
(Caddy, nginx, Cloudflare Tunnel) decides what is reachable from the internet.

## Environments

```
.env.example                       single source of truth for the key list
environments/
  images.manifest.env              third-party images, shared
  staging.manifest.env             application image references
  staging.config.env               runtime configuration
  production.manifest.env
  production.config.env
```

The generated `.env.local`, `.env.staging` and `.env.production` are
**git-ignored**. Secrets never live in the repository; the config files carry
`CHANGE_ME` placeholders and `validate_env.sh` refuses to deploy while any
remain.

## Commands

```bash
make up                  # ENV=local by default
make up ENV=staging
make down ENV=staging
make logs
make ps
make smoke
make validate
```

`make up` and `make down` are also exposed from the repository root.

### Local

```bash
make up
```

Creates `.env.local` from `.env.example` on first run and **builds** the images
from this working copy.

### Staging and production

```bash
./scripts/build_env.sh staging   # assemble .env.staging from environments/
$EDITOR .env.staging             # replace the CHANGE_ME values
./scripts/deploy.sh staging
```

Remote environments never build: they **pull** the image references pinned in
`environments/<env>.manifest.env`. Promoting to production means copying the
verified references from `staging.manifest.env` into `production.manifest.env`.

## Scripts

| Script | Purpose |
| --- | --- |
| `deploy.sh <env> [deploy\|stop]` | Validate, build or pull, start, then smoke-test |
| `smoke.sh <env>` | Nine health checks against a running stack |
| `validate_env.sh <env>` | Every key in `.env.example` present, non-empty, no placeholders |
| `build_env.sh <env>` | Assemble `.env.<env>` from the three manifest files |
| `db.sh <backup\|restore\|list> <env>` | Compressed `pg_dump` / `pg_restore` |

All of them pass `shellcheck`.

## Rollback

Point the manifest at the previous image reference and redeploy:

```bash
$EDITOR environments/production.manifest.env   # restore the previous tag
./scripts/build_env.sh production
./scripts/deploy.sh production
```

Because the manifests are versioned, `git revert` of the promotion commit is
equivalent.

## Backups

```bash
./scripts/db.sh backup production
./scripts/db.sh list production
./scripts/db.sh restore production                       # latest
./scripts/db.sh restore production /path/to/dump.gz      # specific
```

Backups land in `~/backups/midepensa/<env>/` and older than `RETENTION_DAYS`
(7 by default) are pruned on each run. For a nightly dump:

```cron
0 3 * * * cd ~/MiDepensa/infra && ./scripts/db.sh backup production >> ~/backups/midepensa/backup.log 2>&1
```

See [MAINTENANCE.md](MAINTENANCE.md) for day-to-day operations.
