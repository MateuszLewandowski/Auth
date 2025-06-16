LOAD_LOCAL_ENV = set -a; source .env.local; set +a;
LOAD_DEV_ENV = set -a; source .env.dev; set +a;
LOAD_PROD_ENV = set -a; source .env.prod; set +a;

run-local:
	$(LOAD_LOCAL_ENV) docker compose -f docker-compose.local.yml up -d 

stop-local:
	$(LOAD_LOCAL_ENV) docker compose -f docker-compose.local.yml down

run-dev:
	$(LOAD_DEV_ENV) docker compose -f docker-compose.dev.yml up --build

stop-dev:
	$(LOAD_DEV_ENV) docker compose -f docker-compose.dev.yml down	

run-prod:
	$(LOAD_PROD_ENV) docker compose -f docker-compose.prod.yml up --build

stop-prod:
	$(LOAD_PROD_ENV) docker compose -f docker-compose.prod.yml down	

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

test:
	@gotestsum --format testname --format-icons hivis --junitfile unit-tests.xml -- -v ./...

dependencies:
	@go list -m all

docker-logs:
	docker-compose logs -f

sh:
	docker-compose exec -it app sh