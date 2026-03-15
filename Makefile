include .env
export

# Сборка и запуск
run:
	go run cmd/scheduler/main.go

run_bin: build start

build:
	go build -o bin ./cmd/scheduler

start:
	./bin

clean:
	rm bin

# migrations 
MIGRATIONS_DIR=db/migrations

migrate-up:
	migrate -path $(MIGRATIONS_DIR) -database "%(DB_URL)" up

migrate-down:
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" down 1

migrate-drop:
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" drop -f

migrate-force:
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" force $(version)

migrate-version:
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" version

migrate-create:
	@read -p "Migration name: " name; \
	migrate create -seq -ext sql -dir $(MIGRATIONS_DIR) $$name