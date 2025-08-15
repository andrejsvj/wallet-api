-include config.env

.PHONY: local
local:
	env $$(cat config.env | xargs) go run cmd/main.go

.PHONY: build
build:
	go build -o main .

.PHONY: migrate
migrate:
	goose -dir="./migrations" postgres "host=localhost port=$(DB_PORT) user=$(DB_USERNAME) password=$(DB_PASSWORD) database=$(DB_NAME) sslmode=disable" status -v
	goose -dir="./migrations" postgres "host=localhost port=$(DB_PORT) user=$(DB_USERNAME) password=$(DB_PASSWORD) database=$(DB_NAME) sslmode=disable" up -v


.PHONY: new-migration
new-migration:
	goose -dir="./migrations" create $(name) sql

.PHONY: test
test:
	go test -v -short -race ./...



.PHONY: start
start:
	docker-compose -f build/docker-compose.yml --env-file ./config.env up -d

.PHONY: stop
stop:
	docker-compose -f build/docker-compose.yml down

.PHONY: docker-build
docker-build:
	docker-compose -f build/docker-compose.yml build