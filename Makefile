DOCKER_COMPOSE=docker-compose
GOCMD=go
GORUN=$(GOCMD) run
GOTEST=$(GOCMD) test

start-db-dev:
	$(DOCKER_COMPOSE) -f docker/docker-compose.dev.yaml up -d

start-http-dev:
	$(GORUN) ./cmd/http/main.go

start-grpc-dev:
	$(GORUN) ./cmd/grpc/main.go

test:
	$(GOTEST) -v ./... -coverprofile=coverage.out -covermode=atomic