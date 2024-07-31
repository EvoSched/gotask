include .env 	# Include environment variables
.PHONY: 		# List of targets not related to files
.SILENT: 		# Don't show the command executed

build-cli:
	go build -o gt ./cmd/main.go

build-cli-windows:
	GOOS=windows GOARCH=amd64 go build -o gt.exe ./cmd/main.go

build: build-cli build-cli-windows

clean:
	rm -f gt gt.exe

docker-build:
	@docker-compose -f deployments/docker-compose.yml -p gotask --env-file .env up --build

docker-up:
	@docker-compose -f deployments/docker-compose.yml -p gotask --env-file .env up -d

docker-down:
	@docker-compose -f deployments/docker-compose.yml -p gotask --env-file .env down