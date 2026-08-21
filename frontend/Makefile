# MiDepensa Frontend Makefile
#
# Node is optional on the host: when npm is missing every command is executed
# inside a node:22-alpine container instead.

.DEFAULT_GOAL := help

# ─── Toolchain ───────────────────────────────────────────────

NODE_IMAGE ?= node:22-alpine

# Under rootless Docker the container's root already maps to the host user, so
# forcing --user there would break write access to the bind mount.
ifeq ($(shell docker info -f '{{.SecurityOptions}}' 2>/dev/null | grep -c rootless),0)
DOCKER_USER := --user $(shell id -u):$(shell id -g) -e HOME=/tmp
endif

ifeq ($(shell command -v npm 2>/dev/null),)
RUN := docker run --rm -t $(DOCKER_USER) -v $(CURDIR):/app -w /app $(NODE_IMAGE)
else
RUN :=
endif

# ─── .PHONY ─────────────────────────────────────────────────

.PHONY: help deps dev build \
	lint lint-fix fmt fmt-check \
	test test-watch test-ci \
	docker-build docker-start docker-stop \
	clean

# ═══════════════════════════════════════════════════════════════
# Help
# ═══════════════════════════════════════════════════════════════

help:
	@echo "Available commands:"
	@echo ""
	@echo "  Development:"
	@echo "    deps             - Install npm dependencies"
	@echo "    dev              - Dev server with hot reload (http://localhost:4200)"
	@echo "    build            - Production build into www/"
	@echo ""
	@echo "  Code quality:"
	@echo "    lint             - Run ESLint"
	@echo "    lint-fix         - Run ESLint with --fix"
	@echo "    fmt              - Format with Prettier"
	@echo "    fmt-check        - Check formatting without writing"
	@echo ""
	@echo "  Testing:"
	@echo "    test             - Run unit tests (Jest)"
	@echo "    test-watch       - Run tests in watch mode"
	@echo "    test-ci          - Run tests with coverage"
	@echo ""
	@echo "  Docker:"
	@echo "    docker-build     - Build the production image"
	@echo "    docker-start     - Serve the production build with nginx (:4201)"
	@echo "    docker-stop      - Stop the Docker services"
	@echo ""
	@echo "  Maintenance:"
	@echo "    clean            - Remove build artifacts and caches"

# ═══════════════════════════════════════════════════════════════
# Development
# ═══════════════════════════════════════════════════════════════

deps:
	$(RUN) npm install

# Runs through docker compose so the port mapping works either way.
dev:
	docker compose up dev

build: deps
	$(RUN) npx ng build --configuration production

# ═══════════════════════════════════════════════════════════════
# Code quality
# ═══════════════════════════════════════════════════════════════

lint:
	$(RUN) npx eslint src

lint-fix:
	$(RUN) npx eslint src --fix

fmt:
	$(RUN) npx prettier --write "src/**/*.{ts,html,scss}"

fmt-check:
	$(RUN) npx prettier --check "src/**/*.{ts,html,scss}"

# ═══════════════════════════════════════════════════════════════
# Testing
# ═══════════════════════════════════════════════════════════════

test:
	$(RUN) npx jest

test-watch:
	$(RUN) npx jest --watch

test-ci:
	$(RUN) npx jest --ci --coverage

# ═══════════════════════════════════════════════════════════════
# Docker
# ═══════════════════════════════════════════════════════════════

docker-build:
	docker compose build frontend

docker-start: docker-build
	docker compose --profile smoke up -d

docker-stop:
	docker compose --profile smoke down

# ═══════════════════════════════════════════════════════════════
# Maintenance
# ═══════════════════════════════════════════════════════════════

clean:
	rm -rf www dist .angular coverage out-tsc
	@echo "Build artifacts removed"
