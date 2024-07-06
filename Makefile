include .env
.PHONY:
.SILENT:

build-cli:
	go build -o ./.bin/cli ./cmd/cli/main.go

build-tui:
	go build -o ./.bin/tui ./cmd/tui/main.go

build : build-cli build-tui

tui: build-tui
	./.bin/tui