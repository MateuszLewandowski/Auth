# ==============================================================================
# Main commands

run:
	@go run ./cmd/main.go

build:
	@go build ./cmd/main.go

clean:
	@go mod tidy && go fmt ./...

lint:
	@golangci-lint run \
		./config \
		./internal/... \
		./pkg/... \
		./test

up:
	docker compose -f ./docker-compose.yml up -d

test:
	@gotestsum --format testname --format-icons hivis --junitfile unit-tests.xml -- -v ./...

dependencies:
	@go list -m all

down:
	docker compose -f ./docker-compose.yml down

docker-build:
	docker-compose build --no-cache

docker-build-dev:
	docker-compose -f ./docker-compose.yml -f ./docker-compose.override.yml build

docker-up:
	docker-compose up -d
	
docker-down:
	docker-compose down

docker-logs:
	docker-compose logs -f

docker-exec:
	docker-compose exec -it app sh