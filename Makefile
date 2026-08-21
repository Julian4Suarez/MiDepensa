# MiDepensa — root orchestrator
#
# Thin wrapper over backend/Makefile, frontend/Makefile and infra/scripts.
# Every target here delegates; nothing is implemented twice.

.DEFAULT_GOAL := help

ENV ?= local

.PHONY: help \
	up down logs dev \
	test lint fmt \
	smoke deploy \
	db-backup db-restore \
	clean

# ═══════════════════════════════════════════════════════════════
# Help
# ═══════════════════════════════════════════════════════════════

help:
	@echo "MiDepensa — available commands:"
	@echo ""
	@echo "  Stack (Docker):"
	@echo "    up               - Build and start the full stack (postgres + backend + frontend)"
	@echo "    down             - Stop the full stack"
	@echo "    logs             - Tail logs of every service"
	@echo "    smoke            - Health-check a running stack (ENV=local)"
	@echo ""
	@echo "  Development:"
	@echo "    dev              - Postgres in Docker + backend hot-reload + Angular dev server"
	@echo ""
	@echo "  Quality:"
	@echo "    test             - Run backend and frontend tests"
	@echo "    lint             - Lint backend and frontend"
	@echo "    fmt              - Format backend and frontend"
	@echo ""
	@echo "  Operations:"
	@echo "    deploy ENV=staging   - Deploy an environment (local|staging|production)"
	@echo "    db-backup ENV=local  - Dump the database"
	@echo "    db-restore ENV=local - Restore the latest dump"
	@echo ""
	@echo "  Maintenance:"
	@echo "    clean            - Remove build artifacts from both projects"

# ═══════════════════════════════════════════════════════════════
# Stack
# ═══════════════════════════════════════════════════════════════

up:
	@$(MAKE) -C infra up ENV=$(ENV)

down:
	@$(MAKE) -C infra down ENV=$(ENV)

logs:
	@$(MAKE) -C infra logs ENV=$(ENV)

smoke:
	@./infra/scripts/smoke.sh $(ENV)

# ═══════════════════════════════════════════════════════════════
# Development
# ═══════════════════════════════════════════════════════════════

# Starts Postgres, then backend (air) and frontend (ng serve) in parallel.
# Ctrl-C stops both; `make down` cleans up the database container.
dev:
	@$(MAKE) -C backend services-start
	@$(MAKE) -j2 dev-backend dev-frontend

.PHONY: dev-backend dev-frontend
dev-backend:
	@$(MAKE) -C backend dev

dev-frontend:
	@$(MAKE) -C frontend dev

# ═══════════════════════════════════════════════════════════════
# Quality
# ═══════════════════════════════════════════════════════════════

test:
	@$(MAKE) -C backend test
	@$(MAKE) -C frontend test

lint:
	@$(MAKE) -C backend lint
	@$(MAKE) -C frontend lint

fmt:
	@$(MAKE) -C backend fmt
	@$(MAKE) -C frontend fmt

# ═══════════════════════════════════════════════════════════════
# Operations
# ═══════════════════════════════════════════════════════════════

deploy:
	@./infra/scripts/deploy.sh $(ENV)

db-backup:
	@./infra/scripts/db.sh backup $(ENV)

db-restore:
	@./infra/scripts/db.sh restore $(ENV)

# ═══════════════════════════════════════════════════════════════
# Maintenance
# ═══════════════════════════════════════════════════════════════

clean:
	@$(MAKE) -C backend clean
	@$(MAKE) -C frontend clean
