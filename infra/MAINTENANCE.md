# Maintenance

Replace `<env>` with `local`, `staging` or `production`.

## Logs

```bash
docker logs -f midepensa_<env>_backend
docker logs -f midepensa_<env>_frontend
docker compose --env-file .env.<env> logs -f     # everything
```

The backend logs one structured line per request. In `local` the format is text;
elsewhere it is JSON with a `request_id` you can grep for.

## Status

```bash
make ps ENV=<env>
make smoke ENV=<env>

curl -s localhost:8080/healthz    # process is alive
curl -s localhost:8080/readyz     # dependencies answered within 2s
```

`/readyz` returns `503` with `{"dependencies":{"database":"down"}}` when
PostgreSQL is unreachable, which is what a load balancer should act on.

## Database

```bash
# open a shell
docker exec -it midepensa_<env>_postgres psql -U midepensa_user -d midepensa_db

# applied migrations
docker exec midepensa_<env>_postgres psql -U midepensa_user -d midepensa_db \
  -c "SELECT * FROM schema_migrations;"

# how much data there is
docker exec midepensa_<env>_postgres psql -U midepensa_user -d midepensa_db \
  -c "SELECT relname, n_live_tup FROM pg_stat_user_tables ORDER BY n_live_tup DESC;"

# size on disk
docker exec midepensa_<env>_postgres psql -U midepensa_user -d midepensa_db \
  -c "SELECT pg_size_pretty(pg_database_size('midepensa_db'));"
```

Migrations run at startup when `AUTO_MIGRATE=true`. Adding a product to the
catalog means a **new** migration file, never editing `000002_seed_catalog.up.sql`
— golang-migrate records which versions have been applied and will not re-run one.

## Adding a product

1. New migration in `backend/internal/infrastructure/migrations/`, for example
   `000003_add_products.up.sql`, inserting into `products`.
2. Add the code and its emoji codepoint to
   `frontend/scripts/fetch-product-icons.sh` and run it.
3. Redeploy. Existing pantries do **not** get the new item automatically — they
   were seeded at creation time. Add a backfill statement to the same migration
   if you want them updated:

   ```sql
   INSERT INTO pantry_items (pantry_id, product_id, status, pantry_view, category)
   SELECT p.id, pr.id, 'OK', pr.default_view, pr.default_category
   FROM pantries p CROSS JOIN products pr
   ON CONFLICT DO NOTHING;
   ```

## Disk

```bash
docker system df
docker image prune -f      # deploy.sh already does this
```

The only stateful volume is `midepensa_postgres_data`. Back it up before
removing anything.

## Upgrading PostgreSQL

Major versions cannot read an older data directory. Dump, recreate, restore:

```bash
./scripts/db.sh backup <env>
make down ENV=<env>
docker volume rm midepensa_postgres_data
$EDITOR environments/images.manifest.env     # bump POSTGRES_IMAGE
make up ENV=<env>
./scripts/db.sh restore <env>
```
