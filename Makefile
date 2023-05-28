ENV := .env
include $(ENV)

run:
	go run ./cmd/main.go

install:
	go install github.com/matryer/moq@latest

tests:
	go test ./...

subscribe-test:
	curl -X POST http://localhost:9999/subscribe -d email=jon@doe.com

rate:
	curl https://api.coingecko.com/api/v3/exchange_rates | jq '.rates.uah.value'

VERSION := 0.1

docker-build:
	docker build . --tag $(APP_NAME)_v$(VERSION) --file $(DOCKER_FILE) --no-cache

docker-run:
	docker run -it --volume $(DB_PATH):/data $(APP_NAME)_v$(VERSION)

tool-concat:
	./scripts/concat.sh