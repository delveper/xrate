# Genesis Software Engineering School 3.0

## Docs

[gses2swagger.yaml](docs%2Fgses2swagger.yaml)

## Introduction

The application is divided into several key modules as detailed below:

- **cmd**: Contains the application's entry point.
- **data**: Contains file store, or raw data.
- **docs**: Contains documentation files.
- **internal**: Contains the core application logic divided into `rate`, `subscription`, and `transport` packages.
- **scripts**: Contains auxiliary scripts for various tasks.
- **sys**: Contains system-level packages like `env`, `filestore`, and `logger`.

Each module is responsible for a specific function within the application, allowing for clear separation of concerns and
making the codebase easy to manage and navigate.

## Installation and Setup

```shell
make install
```

```shell
make run
```

```shell
make docker-build
 ``` 

```shell
make docker-run
 ```  

## Module Tree

```
📦gentest
 ┣ 📂cmd
 ┃ ┗ 📜main.go
 ┣ 📂data
 ┣ 📂docs
 ┣ 📂internal
 ┃ ┣ 📂rate
 ┃ ┃ ┣ 📜getter_mock_test.go
 ┃ ┃ ┣ 📜handler.go
 ┃ ┃ ┣ 📜handler_test.go
 ┃ ┃ ┣ 📜rate.go
 ┃ ┃ ┗ 📜rate_test.go
 ┃ ┣ 📂subscription
 ┃ ┃ ┣ 📜handler.go
 ┃ ┃ ┣ 📜handler_test.go
 ┃ ┃ ┣ 📜repository.go
 ┃ ┃ ┣ 📜subscriber_mock_test.go
 ┃ ┃ ┗ 📜subscription.go
 ┃ ┗ 📂transport
 ┃   ┣ 📜http.go
 ┃   ┗ 📜middleware.go
 ┣ 📂scripts
 ┣ 📂sys
 ┃ ┣ 📂env
 ┃ ┃ ┣ 📜env.go
 ┃ ┃ ┗ 📜env_test.go
 ┃ ┣ 📂filestore
 ┃ ┃ ┣ 📜filestore.go
 ┃ ┃ ┗ 📜filestore_test.go
 ┃ ┗ 📂logger
 ┃   ┗ 📜logger.go
 ┣ 📜.env
 ┣ 📜.gitignore
 ┣ 📜.golangci.yml
 ┣ 📜Dockerfile
 ┣ 📜go.mod
 ┣ 📜go.sum
 ┣ 📜Makefile
 ┗ 📜README.md
```
