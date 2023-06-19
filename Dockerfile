FROM golang:1.20-alpine as src
WORKDIR /gentest
COPY go.mod  go.sum .env ./
RUN go mod download && go mod verify
ADD . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app ./cmd/main.go

FROM scr AS test
RUN go test ./...

FROM alpine:3.17 as dev
VOLUME ["/data"]
WORKDIR /gentest
COPY --from=src /gentest/app /gentest
COPY --from=src /gentest/.env /gentest
ENTRYPOINT ["/gentest/app"]
