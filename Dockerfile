FROM golang:1.20-alpine as src
WORKDIR /xrate
COPY go.mod  go.sum ./
RUN go mod download && go mod verify
ADD . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main ./app/api/main.go

FROM src AS test
RUN go test ./...

FROM alpine:3.17 as dev
WORKDIR /xrate
COPY --from=src /xrate/main /xrate
ENTRYPOINT ["/xrate/main"]
