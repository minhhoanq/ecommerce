include app.env

proto:
	buf generate api

build:
	go build \
		-o build/app \
		cmd/app/main.go

postgres:
	docker run --name user_service_psql -p 5435:5435 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:17-alpine

createdb:
	docker exec -it user_service_psql createdb --username root --owner root lifeat_user_service_psql

dropdb:
	docker exec -it user_service_psql dropdb lifeat_user_service_psql

migrate-create:
	migrate create -ext sql -dir migrations -seq $(name)

migrateup:
	migrate -path pkg/postgres/migrations -database 'postgresql://${DB_USER}:${DB_PASSWORD}@localhost:${DB_PORT}/${DB_NAME}?sslmode=disable' up

migratedown:
	migrate -path pkg/postgres/migrations -database 'postgresql://${DB_USER}:${DB_PASSWORD}@localhost:${DB_PORT}/${DB_NAME}?sslmode=disable' down

redis:
	docker run --name lifeat-redis -p 6379:6379 -d redis:7-alpine

server:
	go run cmd/app/main.go

evans:
	evans --proto user_service.proto --path api/v1/user_service --host localhost --port 9000 repl

.PHONY: proto build postgres createdb dropdb migrate-create migrateup migratedown redis server proto
