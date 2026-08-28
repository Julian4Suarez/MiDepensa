# Architecture guide

## Domain in one paragraph

A **pantry** is a named collection reachable at `/pantries/{slug}`. Every pantry
holds one **item** per catalog **product**. An item carries three pieces of
state: its place in the shopping workflow (**status**), how often you buy it
(**type**), and what kind of product it is (**category**). The catalog itself
is global, read-only and seeded by a migration.

```
Pantry 1 ──── * PantryItem * ──── 1 Product
"familia-suarez"   PENDING / ESSENTIAL / FRUIT_VEG   "tomato"
```

| Enum | Values |
| --- | --- |
| `ItemStatus` | `DISCARDED`, `PENDING`, `IN_CART`, `ARCHIVED` |
| `ProductType` | `ESSENTIAL`, `SECONDARY` |
| `Category` | `FRUIT_VEG`, `MEAT_FISH`, `DAIRY_EGGS`, `DRY_CANNED`, `DRINKS`, `HOME_CARE` |

Type and category are two independent axes. The first filter also exposes the
archived state; its All option includes only active essential and secondary
products. Type answers *how often do I buy this*, while category mirrors
**supermarket aisles**, because the generated shopping list is grouped by
category and should follow the route you walk through the shop.

The status filter is opened from the toolbar and defaults to `PENDING`. A pending product can be
discarded or placed in the cart; discarded and cart products can return to
pending. The generated shopping list contains every `IN_CART` product and is
independent of the screen filters.
New pantries start with every active product in `PENDING`; the toolbar reset
action restores that state in one bulk request without changing archived items.

Two earlier designs were replaced:

- Four categories (`FRESH`, `PANTRY`, `DRINKS`, `HOME_CARE`) — migration `000003`.
  `PANTRY` collided with the name of the app and `FRESH` held 44% of the
  catalog, which makes it useless as a filter.
- Three *views* (`PRIMARY`, `SECONDARY`, `OTHER`) selected from a side menu —
  migration `000004`. Views and categories both filtered the grid, so having two
  different interactions for the same job was needless. They are now the type
  filter, on the same screen.

## Slugs

`Familia Suárez` becomes `familia-suarez`. The rules live in
[valueobjects/slug.go](../backend/internal/domain/valueobjects/slug.go):
Unicode NFD decomposition drops the accents, every remaining non-alphanumeric
run collapses into one hyphen, and the result is capped at 60 characters.

The same rules are re-implemented in
[slugify.ts](../frontend/src/app/core/utils/slugify.ts) so the home page can
preview the URL as you type. Both have unit tests over the same table of cases;
if they ever drift, one of the two suites fails.

Slugs coming back in from a URL go through `ParseSlug`, which only accepts
already-normalised values. Anything else is a `400` before the database is
touched.

## Backend: hexagonal architecture

Hexagonal architecture (ports and adapters) has one rule: **the core depends on
interfaces, never on technology**. Swapping PostgreSQL for anything else means
writing one new adapter and changing one line in `main.go`.

```
        inbound adapter                core                 outbound adapter
   ┌──────────────────────┐   ┌────────────────────┐   ┌──────────────────────┐
   │  Gin handlers        │──▶│  PantryService     │──▶│  PostgresPantryRepo  │
   │  (HTTP → use case)   │   │  (business rules)  │   │  (port → SQL)        │
   └──────────────────────┘   └────────────────────┘   └──────────────────────┘
        calls a port                                       implements a port
```

| Folder | Role | Depends on |
| --- | --- | --- |
| `internal/domain/entities` | Business objects and their invariants | nothing |
| `internal/domain/valueobjects` | Self-validating values (`Slug`) | nothing |
| `internal/domain/repositories` | Outbound ports for persistence | entities |
| `internal/application/usecases` | Inbound ports consumed by handlers | entities |
| `internal/application/services` | Use case implementations | ports only |
| `internal/application/ports` | Non-persistence outbound ports (`HealthChecker`) | — |
| `internal/infrastructure/adapters/inbound/http` | Gin router, handlers, middleware | usecases |
| `internal/infrastructure/adapters/outbound/persistence` | pgx implementations | repositories |
| `internal/infrastructure/migrations` | Embedded SQL schema and seed | — |
| `internal/config` | Environment variables into a typed struct | — |

Everything is wired once, explicitly, in
[cmd/api/main.go](../backend/cmd/api/main.go). There is no dependency-injection
container: the constructor calls *are* the dependency graph.

### Why services return interfaces they do not name

`services.NewPantryService` returns an unexported struct. `handlers` accepts a
`usecases.PantryService`. The compiler checks the match at the call site in
`main.go`, so the service package never imports the handler package and the
dependency arrow only ever points inwards.

### Error handling

Services return sentinel errors (`services.ErrPantryNotFound`). The HTTP layer
maps them to status codes in exactly one place,
[`registerDomainErrors`](../backend/internal/infrastructure/adapters/inbound/http/router.go).
Handlers call `helpers.Respond(c, err)` and never decide a status code
themselves. Unmapped errors become a logged `500` with a generic body, so
internal details never leak.

| Error | Status | Code |
| --- | --- | --- |
| `ErrPantryNotFound` | 404 | `pantry_not_found` |
| `ErrItemNotFound` | 404 | `item_not_found` |
| `ErrSlugAlreadyExists` | 409 | `slug_already_exists` |
| `ErrInvalidPantryName` | 400 | `invalid_name` |
| `ErrEmptyPatch` | 400 | `empty_update` |
| `ErrInvalidPatch` | 400 | `invalid_update` |

Every error response has the same shape:

```json
{ "error": { "code": "pantry_not_found", "message": "pantry not found" } }
```

## API

Base path `/v1`.

| Method | Path | Purpose |
| --- | --- | --- |
| `GET` | `/healthz` | Liveness — the process is up |
| `GET` | `/readyz` | Readiness — every dependency answered within 2s |
| `GET` | `/v1/catalog` | Products plus the valid enum values |
| `POST` | `/v1/pantries` | `{ "name": "Familia Suárez" }` → 201 with the slug |
| `GET` | `/v1/pantries/{slug}` | Pantry with all items |
| `PATCH` | `/v1/pantries/{slug}/items/{productId}` | `{ status?, type?, category? }` |
| `POST` | `/v1/pantries/{slug}/items/reset` | Reset every active item to `PENDING` |

Creating a pantry inserts the pantry row and all 57 item rows in one
transaction; a slug collision rolls the whole thing back and returns `409`. The
home page turns that `409` into navigation: the pantry the user asked for
already exists, so it opens it instead of reporting an error.

The `PATCH` is a real partial update. `nil` fields become SQL `NULL` and the
statement uses `COALESCE($3, status)`, so omitted fields keep their value. A CTE
returns the updated row joined with its product, so one round trip is enough.

## Middleware chain

Applied in order in `SetupRouter`:

1. `RequestID` — reuses or generates `X-Request-Id`
2. `Recovery` — a panic becomes a logged `500` instead of a dead process
3. `Logger` — one structured line per request
4. `SecurityHeaders` — `nosniff`, `DENY`, `no-referrer`
5. `CORS` — explicit origin allowlist
6. `BodyLimit` — 8 KiB, well above the largest legitimate payload

## Database

PostgreSQL 16. Three tables, defined in
[000001_init.up.sql](../backend/internal/infrastructure/migrations/000001_init.up.sql):

- `pantries` — `slug` is `UNIQUE`, which is what makes the `409` race-proof
- `products` — the seeded catalog
- `pantry_items` — composite primary key `(pantry_id, product_id)`,
  `ON DELETE CASCADE` from the pantry

Every enum column has a `CHECK` constraint, so invalid data cannot exist even if
a bug slips past the Go validation. Renaming an enum therefore means dropping
the constraint, rewriting the values and adding it back — see migrations `000003`
and `000004`.

Migrations are embedded with `//go:embed *.sql`, so the binary carries its own
schema and the container image needs no extra files.

## Security

Deliberately proportionate to a family pantry tracker with no accounts:

- parameterised SQL everywhere (no string concatenation)
- slug and UUID validation before any query
- 8 KiB request body cap
- CORS allowlist, `*` only in local development
- non-root users in both container images
- ports published on `127.0.0.1` outside local development
- no secrets reach the browser: the frontend only receives `BACKEND_URL`

There is no authentication, by design. Anyone with the link can edit the pantry.

## Testing strategy

| Level | Location | Needs | What it proves |
| --- | --- | --- | --- |
| Unit | `backend/test/unit` | nothing | Slug rules, entity invariants, service branching |
| Integration | `backend/test/integration` | PostgreSQL | SQL, constraints, partial updates, scanning |
| Unit (frontend) | `frontend/src/**/*.spec.ts` | nothing | Slugify, shopping list, sorting, store filtering and rollback |

Mocks are hand-written structs with function fields — no mocking framework.
Integration tests are behind the `integration` build tag so `go test ./...`
stays fast and dependency-free.

The Jest transformer strips types without checking them, so `make -C frontend
lint` runs `tsc --noEmit` over `tsconfig.spec.json` first: without it a wrong
enum literal in a spec only shows up as a confusing failed assertion.
