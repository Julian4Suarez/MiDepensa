# MiDepensa

A small, anonymous pantry-stock tracker. Create a pantry, mark each product as
**enough**, **low** or **out**, and turn it into a shopping list in one tap.

No sign-up, no accounts: whoever has the link has the pantry.

```
Home  →  "Familia Suárez"  →  /pantries/familia-suarez
```

## Stack

| Layer | Technology |
| --- | --- |
| Frontend | Ionic 8 + Angular 20 (standalone components, signals), Capacitor-ready |
| Backend | Go 1.26 + Gin, hexagonal architecture |
| Database | PostgreSQL 16, migrations embedded in the binary |
| Infrastructure | Docker Compose + Make, one command per environment |

## Quickstart

Requires Docker and Make. Node and Go are optional — everything runs in containers.

```bash
make up      # build and start postgres + backend + frontend, then health-check
make down    # stop everything
```

Then open <http://localhost:4200>.

| Service | URL |
| --- | --- |
| Frontend | <http://localhost:4200> |
| API | <http://localhost:8080/v1> |
| Liveness / readiness | <http://localhost:8080/healthz>, <http://localhost:8080/readyz> |

## Development

```bash
make dev     # postgres in Docker + Go hot reload + Angular dev server
make test    # backend and frontend tests
make lint    # go vet + golangci-lint + ESLint
make fmt     # gofmt + Prettier
make smoke   # health-check a running stack
```

Each project also has its own Makefile — run `make help` in `backend/`,
`frontend/` or `infra/` for the full list.

## Repository layout

```
backend/    Go API (hexagonal: domain / application / infrastructure)
frontend/   Ionic + Angular app
infra/      Docker Compose, environment manifests, deploy and smoke scripts
docs/       Setup, architecture and Ionic guides
```

## Documentation

| Document | What it covers |
| --- | --- |
| [docs/GUIDE_SETUP.md](docs/GUIDE_SETUP.md) | Prerequisites, environment variables, daily workflow |
| [docs/GUIDE_ARCHITECTURE.md](docs/GUIDE_ARCHITECTURE.md) | Domain model, API, hexagonal layers |
| [docs/GUIDE_IONIC_FRONTEND.md](docs/GUIDE_IONIC_FRONTEND.md) | How the Ionic/Angular app works, top to bottom |
| [infra/README.md](infra/README.md) | Environments, deploys, rollbacks |
| [CONTRIBUTING.md](CONTRIBUTING.md) | Branch and commit conventions |

## Credits

Product icons come from [OpenMoji](https://openmoji.org) (CC BY-SA 4.0).
See [frontend/ATTRIBUTION.md](frontend/ATTRIBUTION.md).

