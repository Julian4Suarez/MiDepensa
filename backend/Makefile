# MiDepensa Backend Makefile

.DEFAULT_GOAL := help

# ─── Variables ───────────────────────────────────────────────

GOBIN := $(shell go env GOPATH)/bin
CI ?= 0
COVERAGE ?= 0

TEST_FLAGS := -count=1
ifeq ($(CI),1)
TEST_FLAGS += -race
endif
ifeq ($(COVERAGE),1)
TEST_FLAGS += -coverprofile=coverage.out -covermode=atomic
endif

# ─── .PHONY ─────────────────────────────────────────────────

.PHONY: help deps dev dev-stop app-build app-start \
	services-start services-stop \
	test test-unit test-integration test-clean \
	fmt lint \
	docker-build docker-start docker-stop \
	validate-branch changelog clean

# ═══════════════════════════════════════════════════════════════
# Help
# ═══════════════════════════════════════════════════════════════

help:
	@echo "Available commands:"
	@echo ""
	@echo "  Development:"
	@echo "    deps             - Install dev tools (air, golangci-lint)"
	@echo "    dev              - Start postgres and run the API with hot reload"
	@echo "    dev-stop         - Stop the hot-reload process"
	@echo "    app-build        - Compile the binary into bin/api"
	@echo "    app-start        - Build and run the API in the background"
	@echo ""
	@echo "  Services:"
	@echo "    services-start   - Start postgres in Docker"
	@echo "    services-stop    - Stop postgres"
	@echo ""
	@echo "  Testing:"
	@echo "    test             - Run every test"
	@echo "    test-unit        - Run unit tests (no dependencies)"
	@echo "    test-integration - Run integration tests (needs postgres)"
	@echo "    test-clean       - Clear the Go test cache"
	@echo "    Flags: CI=1 (race detector) COVERAGE=1 (coverage profile)"
	@echo ""
	@echo "  Code quality:"
	@echo "    fmt              - Format with gofmt"
	@echo "    lint             - Run go vet and golangci-lint if installed"
	@echo ""
	@echo "  Docker:"
	@echo "    docker-build     - Build the production image"
	@echo "    docker-start     - Run postgres + API in Docker"
	@echo "    docker-stop      - Stop the Docker services"
	@echo ""
	@echo "  CI:"
	@echo "    validate-branch  - Validate the branch name convention"
	@echo "    changelog        - Regenerate CHANGELOG.md with git-cliff"
	@echo ""
	@echo "  Maintenance:"
	@echo "    clean            - Remove build and coverage artifacts"

# ═══════════════════════════════════════════════════════════════
# Development
# ═══════════════════════════════════════════════════════════════

deps:
	go install github.com/air-verse/air@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@echo "Installed into $(GOBIN)"

# Hot reload; falls back to a plain run when air is not installed.
dev: services-start
	@if [ -x "$(GOBIN)/air" ]; then \
		$(GOBIN)/air; \
	else \
		echo "air not found — running without hot reload (make deps to install)"; \
		go run ./cmd/api; \
	fi

dev-stop:
	@pkill -f '[a]ir' || true

app-build:
	go build -ldflags="-X midepensa/internal/config.Commit=$$(git rev-parse --short HEAD 2>/dev/null || echo dev)" \
		-o bin/api ./cmd/api

app-start: app-build
	./bin/api &

# ═══════════════════════════════════════════════════════════════
# Services
# ═══════════════════════════════════════════════════════════════

services-start:
	docker compose up -d --wait postgres

services-stop:
	docker compose down

# ═══════════════════════════════════════════════════════════════
# Testing
# ═══════════════════════════════════════════════════════════════

test: test-unit test-integration

test-unit:
	go test $(TEST_FLAGS) ./test/unit/...

# Starts postgres first so `make test` works from a clean checkout.
test-integration: services-start
	go test $(TEST_FLAGS) -tags=integration ./test/integration/...

test-clean:
	go clean -testcache

# ═══════════════════════════════════════════════════════════════
# Code quality
# ═══════════════════════════════════════════════════════════════

fmt:
	gofmt -w .

lint:
	go vet ./...
	@if [ -x "$(GOBIN)/golangci-lint" ]; then \
		$(GOBIN)/golangci-lint run; \
	else \
		echo "golangci-lint not installed — go vet only (make deps to install)"; \
	fi

# ═══════════════════════════════════════════════════════════════
# Docker
# ═══════════════════════════════════════════════════════════════

docker-build:
	docker compose build \
		--build-arg APP_COMMIT=$$(git rev-parse --short HEAD 2>/dev/null || echo dev) \
		--build-arg APP_BUILD_TIME=$$(date -u +%Y-%m-%dT%H:%M:%SZ) \
		app

docker-start: docker-build
	docker compose --profile smoke up -d --wait

docker-stop:
	docker compose --profile smoke down

# ═══════════════════════════════════════════════════════════════
# CI
# ═══════════════════════════════════════════════════════════════

validate-branch:
	@BRANCH=$${BRANCH:-$$(git rev-parse --abbrev-ref HEAD)}; \
	echo "Validating branch: $$BRANCH"; \
	if [ "$$BRANCH" = "main" ]; then \
		echo "Protected branch ($$BRANCH)"; exit 0; \
	fi; \
	if echo "$$BRANCH" | grep -qE '^(feat|fix|chore)/[a-z0-9._-]+$$'; then \
		echo "Branch name is valid"; \
	else \
		echo "Invalid branch name: $$BRANCH"; \
		echo "  Must match: feat|fix|chore/<kebab-case-description>"; \
		exit 1; \
	fi

changelog:
	git-cliff --output CHANGELOG.md

# ═══════════════════════════════════════════════════════════════
# Maintenance
# ═══════════════════════════════════════════════════════════════

clean:
	rm -rf bin tmp coverage.out coverage.html
	@echo "Build artifacts removed"
