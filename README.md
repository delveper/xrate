t # Genesis Software Engineering School 3.0

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
    App & Handlers & SubscriptionAdapters & RateAdapters -->|uses| Web
    App -->|binds| RateService & SubscriptionService & Infrastructure & RateAdapters & SubscriptionAdapters
    Domain ==> Handlers
    RateAdapters -.->|impl| ExchangeRateProvider
    SubscriptionAdapters -.->|impl| EmailSender
    SubscriptionService -.->|impl| SubscriptionServiceInterface
    RateService -.->|impl| RateServiceInterface
    Client[Client] -->|interacts| HTTP
    main-->|serves| HTTP
    subgraph Transport
        subgraph HTTP
            App((APP)) -->|binds| RateHandlers[Rate Handlers]
            App -->|binds| SubscriptionHandlers[Subscription Handlers]
            subgraph Handlers
                RateHandlers[/Rate Handlers/] -->|uses| RateServiceInterface{{RateService}}
                SubscriptionHandlers[/Subscription Handlers/] -->|uses| SubscriptionServiceInterface{{SubscriptionService}}
            end
        end
    end
    subgraph RateAdapters
        A
        B
        C
        D
    end
    subgraph SubscriptionAdapters
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
                Message(Message)
                Topic(Topic)
            end
            SubscriptionService((SERVICE)) --> SubscriptionCore
            SubscriptionService -->|uses| Repository{{SubscriberRepository}}
            SubscriptionService -->|uses| EmailSender{{EmailSender}}
            SubscriptionService -->|uses| SubscriptionEvent
        end
    end
    subgraph Infrastructure
        subgraph Env
        end
        Repository -.->|impl| FileStore[(File Store)]
        subgraph Event
            SubscriptionEvent{Event} -->|uses| EventBus((Event Bus))
            RateEvent{Event} -->|uses| EventBus((Event Bus))
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
    
    class SubscriptionHandler{
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


## Sending Emails
```mermaid
sequenceDiagram
    participant Client as Client
    participant Service as Service
    participant SMTPClient as SMTPClient
    Client ->> Service: Send(Message)
    Service ->> SMTPClient: Connect()
    SMTPClient -->> Service: Connection established
    Service ->> SMTPClient: Authenticate()
    SMTPClient -->> Service: Authentication successful
    Service ->> SMTPClient: Mail()
    SMTPClient -->> Service: Mail command successful
    Service ->> SMTPClient: Recipient(Message)
    SMTPClient -->> Service: Recipient command successful
    Service ->> SMTPClient: Data(Message)
    SMTPClient -->> Service: Data command successful
    Service ->> SMTPClient: Quit()
    SMTPClient -->> Service: Quit command successful
```