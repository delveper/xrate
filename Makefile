ENV := .env
include $(ENV)

VERSION := 1.0.0
DOCKER_CONFIG_FLAGS := --file $(DOCKER_COMPOSE_FILE) --env-file $(ENV) --log-level $(LOG_LEVEL)

run:
	go run ./app/api/main.go

install:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin v1.53.3

generate:
	go generate ./...

tests:
	go test -count=1 -v ./...

tests-int:
	INTEGRATION_RUN=true go test -count=1 -v ./...

tests-e2e:
	docker-compose ${DOCKER_CONFIG_FLAGS} up --abort-on-container-exit

consume-log:
	go run ./app/logcs/main.go

lint:
	golangci-lint run -c .golangci.yml

docker-up:
	docker-compose ${DOCKER_CONFIG_FLAGS} up --detach

docker-down:
	docker-compose ${DOCKER_CONFIG_FLAGS} down --remove-orphans

docker-build:
	docker-compose ${DOCKER_CONFIG_FLAGS} build --no-cache
