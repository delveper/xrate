# Genesis Software Engineering School 3.0

## Doc

[openapi.yaml](doc%2Fopenapi.yaml)

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

--- TODO: Update

```
📦xrate
 ┣ 📂.github
 ┃ ┗ 📂workflows
 ┃   ┣ 📜go.yml
 ┃   ┗ 📜golangci.yml
 ┣ 📂api
 ┃ ┣ 📜api.go
 ┃ ┣ 📜config.go
 ┃ ┗ 📜routes.go
 ┣ 📂cmd
 ┃ ┗ 📜main.go
 ┣ 📂doc
 ┃ ┗ 📜openapi.yaml
 ┣ 📂internal
 ┃ ┣ 📂rate
 ┃ ┃ ┣ 📜config.go
 ┃ ┃ ┣ 📂curxrt
 ┃ ┃ ┃ ┣ 📜alphavantage.go
 ┃ ┃ ┃ ┣ 📜coinapi.go
 ┃ ┃ ┃ ┣ 📜coinyep.go
 ┃ ┃ ┃ ┣ 📜curxrt.go
 ┃ ┃ ┃ ┣ 📜ninjas.go
 ┃ ┃ ┃ ┗ 📜xratehost.go
 ┃ ┃ ┣ 📜event.go
 ┃ ┃ ┣ 📜handler.go
 ┃ ┃ ┗ 📜rate.go
 ┃ ┗ 📂subs
 ┃   ┣ 📜config.go
 ┃   ┣ 📜event.go
 ┃   ┣ 📜handler.go
 ┃   ┣ 📜repo.go
 ┃   ┣ 📜repo_test.go
 ┃   ┣ 📜sender.go
 ┃   ┗ 📜subs.go
 ┣ 📂log
 ┃ ┗ 📜sys.log
 ┣ 📂sys
 ┃ ┣ 📂env
 ┃ ┃ ┣ 📜env.go
 ┃ ┃ ┗ 📜env_test.go
 ┃ ┣ 📂event
 ┃ ┃ ┗ 📜event.go
 ┃ ┣ 📂filestore
 ┃ ┃ ┣ 📜filestore.go
 ┃ ┃ ┗ 📜filestore_test.go
 ┃ ┣ 📂logger
 ┃ ┃ ┗ 📜logger.go
 ┃ ┗ 📂web
 ┃   ┣ 📜errors.go
 ┃   ┣ 📜middlewares.go
 ┃   ┣ 📜middlewares_test.go
 ┃   ┣ 📜params.go
 ┃   ┣ 📜request.go
 ┃   ┣ 📜respond.go
 ┃   ┗ 📜web.go
 ┣ 📂test
 ┃ ┣ 📂mock
 ┃ ┃ ┣ 📜email_repository.go
 ┃ ┃ ┣ 📜email_sender.go
 ┃ ┃ ┣ 📜getter.go
 ┃ ┃ ┗ 📜subscriber.go
 ┃ ┣ 📜Dockerfile
 ┃ ┗ 📜postman.json
 ┣ 📜.gitignore
 ┣ 📜.golangci.yml
 ┣ 📜Dockerfile
 ┣ 📜Makefile
 ┣ 📜README.md
 ┣ 📜docker-compose.yml
 ┣ 📜go.mod
 ┗ 📜go.sum

```

## Architecture

```mermaid
graph TB
    main((main)) ==> App
    main ==> Env
    main & EventBus & Web & App ==> Logger>Logger]
    App & Handlers & NotificationAdapters & RateAdapters -->|uses| Web
    App -->|binds| RateService & SubscriptionService & NotificationService & Infrastructure & RateAdapters & NotificationAdapters
    Domain ==> Handlers
    RateAdapters -.->|impl| ExchangeRateProvider
    NotificationAdapters -.->|impl| EmailSender
    SubscriptionService -.->|impl| SubscriptionServiceInterface
    RateService -.->|impl| RateServiceInterface
    NotificationService -.->|impl|NotificationServiceInterface
    Client[Client] -->|interacts| HTTP
    main -->|serves| HTTP
    subgraph Transport
        subgraph HTTP
            App((APP)) -->|binds| RateHandlers[Rate Handlers]
            App -->|binds| SubscriptionHandlers[Subscription Handlers]
            subgraph Handlers
                RateHandlers[/Rate Handlers/] -->|uses| RateServiceInterface{{RateService}}
                SubscriptionHandlers[/Subscription Handlers/] -->|uses| SubscriptionServiceInterface{{SubscriptionService}}
                NotificationHandlers[/Notification Handlers/] -->|uses| NotificationServiceInterface{{NotificationService}}
            end
        end
    end
    subgraph RateAdapters
        A
        B
        C
        D
    end
    subgraph NotificationAdapters
        EmailClient
    end
    subgraph Domain
        subgraph Rate
            subgraph RateCore
                ExchangeRate(ExchangeRate)
            end
            RateService((SERVICE)) --> ExchangeRate
            RateService -->|uses| RateEvent
            RateService -->|uses| ExchangeRateProvider{{ExchangeRateProvider}}
        end
        subgraph Subscription
            subgraph SubscriptionCore
                Subscriber{Subscriber}
                Topic(Topic)
            end
            SubscriptionService((SERVICE)) --> SubscriptionCore
            SubscriptionService -->|uses| Repository{{SubscriberRepository}}
            SubscriptionService -->|uses| SubscriptionEvent
        end
        subgraph Notification
            subgraph NotificationCore
                Message(Message)
                Topic(Topic)
            end
            NotificationService((SERVICE)) --> NotificationCore
            NotificationService -->|uses| MessageCreator{{MessageCreator}}
            NotificationService -->|uses| EmailSender{{Sender}}
            NotificationService -->|uses| NotificationEvent
        end
    end
    subgraph Infrastructure
        subgraph Env
        end
        Repository -.->|impl| FileStore[(File Store)]
        subgraph Event
            SubscriptionEvent{Event} -->|uses| EventBus((Event Bus))
            RateEvent{Event} -->|uses| EventBus((Event Bus))
            NotificationEvent{Event} -->|uses| EventBus((Event Bus))
        end
        subgraph Web
            Middleware
            Tooling
        end
    end
```

## Entities

--TODO: Finish

```mermaid
classDiagram
    class App {
        <<struct>>
        sig chan os.Signal
        log *logger.Logger
        web *web.Web
    }
    class Route {
        <<type>>
    }
    class ConfigAggregate {
        <<struct>>
        Api Config
        Rate rate.Config
        Subscription subs.Config
    }
    class Config {
        <<struct>>
        Name string
        Path string
        Version string
        Origin string
    }
    class RateHandler {
        <<struct>>
        rate ExchangeRateService
    }

    class SubscriptionHandler {
        subs SubscriptionService
    }

    class ExchangeRateService {
        <<interface>>
        GetExchangeRate(ctx context.Context, currency CurrencyPair) (*ExchangeRate, error)
    }
    class Web {
        <<struct>>
        mux *httprouter.Router
        mws []Middleware
        sig chan os.Signal
    }
    class Middleware {
        <<type>>
    }
    class SubscriptionService {
        <<interface>>
        Subscribe(context.Context, Subscriber) error
        SendEmails(context.Context) error
    }
    class Response {
        <<struct>>
        Message string
    }
    class Subscriber {
        <<struct>>
        Address *mail.Address
        Topic Topic
    }
    class RateConfig {
        <<struct>>
        Provider struct
        Client struct
    }
    class ProviderConfig {
        <<struct>>
        Name string
        Endpoint string
        Header string
        Key string
    }
    class SubsConfig {
        <<struct>>
        Sender SenderConfig
        Repo RepoConfig
    }
    class SenderConfig {
        <<struct>>
        Address string
        Key string
    }
    class RepoConfig {
        <<struct>>
        Data string
    }
    class Storer {
        <<interface>>
        Store(Subscriber) error
        FetchAll() ([]Subscriber, error)
    }
    class Repo {
        <<struct>>
        Storer
    }
    class Logger {
        <<struct>>
        *zap.SugaredLogger
    }


    App o-- Route
    App --> ConfigAggregate
    App --> Web
    App --> Logger
    ConfigAggregate o-- Config
    ConfigAggregate o-- RateConfig
    ConfigAggregate o-- SubsConfig
    Handler o-- ExchangeRateService
    Web -- Middleware
    SubscriptionService -- Subscriber
    SubscriptionService -- Response
    RateConfig o-- ProviderConfig
    SubsConfig o-- SenderConfig
    SubsConfig o-- RepoConfig
    Repo o-- Storer
```