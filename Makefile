.PHONY: build run test lint docker-build docker-up docker-down migrate

build:
	go build -o bin/server ./cmd/server

run:
	go run ./cmd/server

test:
	go test -v -race -coverprofile=coverage.out ./...

lint:
	golangci-lint run

docker-build:
	docker-compose build

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

migrate:
	go run ./cmd/migrate

seed:
	go run ./cmd/seed

swagger:
	@echo "Generating Swagger documentation..."
	@which swag > /dev/null || (echo "Installing swag..." && go install github.com/swaggo/swag/cmd/swag@latest)
	swag init -g cmd/server/main.go -o docs

clean:
	rm -rf bin/
	rm -f coverage.out
	rm -rf docs/

