FROM golang:1.20-alpine as src
WORKDIR /gensch
COPY go.mod  go.sum .env ./
RUN go mod download && go mod verify
ADD . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app ./cmd/main.go

FROM src AS test
RUN go test ./...

FROM alpine:3.17 as dev
WORKDIR /gensch
COPY --from=src /xrate/app /xrate
COPY --from=src /xrate/.env /xrate
ENTRYPOINT ["/xrate/app"]
