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