# Variables
include .env.dev
export

DOCKER_COMPOSE = docker-compose --file build/docker-compose.base.yml --file build/docker-compose.dev.yml --project-directory . --project-name ${PROJECT_NAME}
DOCKER = docker

# Run
.PHONY: setup dev build
setup: pg pg-migrate build-dev-img

dev:
	@${DOCKER_COMPOSE} run --service-ports --rm api sh -c 'go run ./cmd/serverd'

build:
	@${DOCKER_COMPOSE} run --rm api sh -c 'go build -o ./output/serverd ./cmd/serverd'

run:
	@${DOCKER_COMPOSE} run --service-ports --rm api sh -c './output/serverd'

# Helper
.PHONY: build-dev-img generate docker-compose-config teardown
build-dev-img:
	@${DOCKER} build -f build/api-dev.Dockerfile -t ${PROJECT_NAME}-go-dev:latest . --build-arg PROJECT_NAME=${PROJECT_NAME}

generate:
	@${DOCKER_COMPOSE} run --rm api sh -c 'go generate ./...'

docker-compose-config:
	@${DOCKER_COMPOSE} config

teardown:
	@${DOCKER_COMPOSE} down

# Database
.PHONY: pg pg-migrate pg-drop
pg:
	@${DOCKER_COMPOSE} up -d pg

pg-migrate:
	@${DOCKER_COMPOSE} run --rm pg-migrate sh -c 'migrate -path /migrations -database "$$PG_URL" up'

pg-drop:
	@${DOCKER_COMPOSE} run --rm pg-migrate sh -c 'migrate -path /migrations -database "$$PG_URL" drop'

# Collector
.PHONY: metrics collector jaeger prometheus grafana
metrics: jaeger prometheus collector grafana

collector:
	@${DOCKER_COMPOSE} up -d collector

jaeger:
	@${DOCKER_COMPOSE} up -d jaeger

prometheus:
	@${DOCKER_COMPOSE} up -d prometheus

grafana:
	@${DOCKER_COMPOSE} up -d grafana
