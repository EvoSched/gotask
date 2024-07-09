include .env
.PHONY:
.SILENT:

build-cli:
	go build -o task ./cmd/cli/main.go

build-tui:
	go build -o tui ./cmd/tui/main.go

build : build-cli build-tui

tui: build-tui
	./.bin/tui

docker-build:
	@docker-compose -f deployments/docker-compose.yml -p gotask --env-file .env up --build

docker-up:
	@docker-compose -f deployments/docker-compose.yml -p gotask --env-file .env up -d

docker-down:
	@docker-compose -f deployments/docker-compose.yml -p gotask --env-file .env down