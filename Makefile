ENV := .env
include $(ENV)

run:
	go run ./cmd/main.go

install:
	go install github.com/matryer/moq@latest
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.53.3

generate:
	go generate ./...

test:
	go test -count=1 ./...

lint:
	golangci-lint run -c .golangci.yml

rate-service:
	curl https://api.coingecko.com/api/v3/exchange_rates | jq '.rates.uah.value'

VERSION := 1.0.0

docker-build:
	docker build . --tag $(APP_NAME)_v$(VERSION) --file $(DOCKER_FILE) --no-cache

docker-run:
	docker run -it --volume $(DB_PATH):/data $(APP_NAME)_v$(VERSION)