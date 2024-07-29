include .env 	# Include environment variables
.PHONY: 		# List of targets not related to files
.SILENT: 		# Don't show the command executed

build-cli:
	go build -o gt ./cmd/main.go

build : build-cli

cli: build-cli
	./gt

docker-build:
	@docker-compose -f deployments/docker-compose.yml -p gotask --env-file .env up --build

docker-up:
	@docker-compose -f deployments/docker-compose.yml -p gotask --env-file .env up -d

docker-down:
	@docker-compose -f deployments/docker-compose.yml -p gotask --env-file .env down