ifneq ("$(wildcard ./.env)","")
    include .env
	export
endif

DB_DSN="postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSL_MODE)"

dev:
	@docker compose -f compose.dev.yml up --build

prod:
	@docker compose -f compose.yml up --build

go.mod.tidy:
	@go mod tidy -e

migrate.create:
	@migrate create -ext sql -dir ./migrations -seq $(name)

migrate.version:
	@migrate -database $(DB_DSN) -path=./migrations version

migrate.all.up:
	@migrate -database $(DB_DSN) -path ./migrations up

migrate.all.down:
	@migrate -database $(DB_DSN) -path ./migrations down

migrate.goto:
	@migrate -database $(DB_DSN) -path ./migrations goto $(V)