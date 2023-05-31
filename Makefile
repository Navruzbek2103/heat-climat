POSTGRES_HOST=localhost
POSTGRES_PORT=5432
POSTGRES_USER=developer
POSTGRES_PASSWORD=2002
POSTGRES_DATABASE=climatedb

DB_URL="postgresql://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DATABASE)?sslmode=disable"

migrate_file:
	migrate create -ext sql -dir migrations -seq create_tablename_table

run:
	go run cmd/main.go

migrate_up:
	migrate -path migrations -database "$(DB_URL)" -verbose up

migrate_down:
	migrate -path migrations -database "$(DB_URL)" -verbose down 
local-up:
	docker compose --env-file ./.env.docker up -d

migrate_forse:
	migrate -path migrations -database "$(DB_URL)" -verbose forse

swag:
	swag init -g api/router.go -o api/docs