ENV := .env
include $(ENV)

run:
	go run ./cmd/main.go

install:
	go install github.com/matryer/moq@latest

tests:
	go test ./...

subscribe-test:
	curl -X POST http://localhost:9999/api/subscribe -d email=jon@doe.com

rate-test:
	curl http://localhost:9999/api/rate

rate-service:
	curl https://api.coingecko.com/api/v3/exchange_rates | jq '.rates.uah.value'

VERSION := 1.0.0

docker-build:
	docker build . --tag $(APP_NAME)_v$(VERSION) --file $(DOCKER_FILE) --no-cache

docker-run:
	docker run -it --volume $(DB_PATH):/data $(APP_NAME)_v$(VERSION)

concat:
	./scripts/concat.sh