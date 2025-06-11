build:
    go build -o bin/api ./cmd/api

run: build
    source .envrc.local && ./bin/api

test:
    go test ./...

test-race:
    go test -race ./...

dev:
    source .envrc.local && air --build.cmd "go build -o bin/api cmd/api/main.go" --build.bin "./bin/api"

db-up:
    docker-compose up -d postgres

db-down:
    docker-compose down

db-logs:
    docker-compose logs -f postgres

help:
    @just --list 