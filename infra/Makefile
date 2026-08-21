# MiDepensa Infrastructure Makefile
#
# Thin wrapper over scripts/. The root Makefile delegates here.

.DEFAULT_GOAL := help

ENV ?= local

.PHONY: help up down restart logs ps smoke validate deploy

help:
	@echo "Infrastructure commands (override the environment with ENV=...):"
	@echo ""
	@echo "  up        - Build and start the stack        (ENV=$(ENV))"
	@echo "  down      - Stop the stack                   (ENV=$(ENV))"
	@echo "  restart   - Stop and start again             (ENV=$(ENV))"
	@echo "  logs      - Tail logs of every service       (ENV=$(ENV))"
	@echo "  ps        - Show container status            (ENV=$(ENV))"
	@echo "  smoke     - Health-check a running stack     (ENV=$(ENV))"
	@echo "  validate  - Validate .env.$(ENV) against .env.example"
	@echo "  deploy    - Alias of up"

up deploy:
	@./scripts/deploy.sh $(ENV)

down:
	@./scripts/deploy.sh $(ENV) stop

restart: down up

logs:
	@docker compose --env-file .env.$(ENV) logs -f

ps:
	@docker compose --env-file .env.$(ENV) ps

smoke:
	@./scripts/smoke.sh $(ENV)

validate:
	@./scripts/validate_env.sh $(ENV)
