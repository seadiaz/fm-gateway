build:
    go build -o bin/api ./cmd/api

run: build
    source .envrc.local
    ./bin/api

test:
    go test ./...

test-race:
    go test -race ./...

dev:
    air --build.cmd "go build -o bin/api cmd/api/main.go" --build.bin "./bin/api"

help:
    @just --list 