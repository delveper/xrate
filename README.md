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

## Architecture

```mermaid
graph TD
    main --> App
    main --> Env
    main --> Logger>Logger]
    EventBus --> Logger
    Web --> Logger
    App -->|uses| Web

    subgraph "Transport"
        App((APP)) -->|binds| RateHandlers[Rate Handlers]
        App((APP)) -->|binds| SubscriptionHandlers[Subscription Handlers]
    end
    subgraph "Domain"
        subgraph "Rate"
            ExchangeRate(ExchangeRate) --> RateService((SERVICE))
            RateHandlers -->|uses| RateServiceInterface{{RateService}}
            RateService((SERVICE)) -->|impl| RateServiceInterface{{RateService}}
            RateService((SERVICE)) -->|uses| ExchangeRateServiceInterface{{ExchangeRateService}}
            RateService((SERVICE)) -->|uses| RateEvent
            RateService((SERVICE)) -->|uses| ExchangeRateProvider{{ExchangeRateProvider}}
        end
        subgraph "Subscription"
            SubscriberEntity{Subscriber} --> SubscriptionService((SERVICE))
            MessageEntity(Message) --> SubscriptionService((SERVICE))
            TopicEntity(Topic) --> SubscriptionService((SERVICE))
            SubscriptionHandlers -->|uses| SubscriptionServiceInterface{{SubscriptionService}}
            SubscriptionService((SERVICE)) -->|implements| SubscriptionServiceInterface
            SubscriptionService -->|uses| Repository{{SubscriberRepository}}
            SubscriptionService -->|uses| EmailSender{{EmailSender}}
            SubscriptionService -->|uses| SubscriptionEvent

        end
    end

    ExchangeRateProvider{{ExchangeRateProvider}} -->|implements| RateAdapters
    subgraph RateAdapters
        A
        B
        C
        D
    end

    EmailSender -->|implements| SubscriptionAdapters
    subgraph SubscriptionAdapters
        EmailClient
    end
    RateAdapters -->|uses| Web[[Web]]
    SubscriptionAdapters -->|uses| Web[[Web]]

    subgraph "Infrastructure"
        Repository -->|Implement| FileStore[(File Store)]
        App((Controller)) -->|uses| Web[[Web]]
        RateHandlers[/Rate Handler/] -->|uses| Web[[Web]]
        SubscriptionHandlers[/Subscription Handler/] -->| uses| Web[[Web]]
        subgraph Web
            Middleware
            Tooling
        end
        subgraph Env
        end
    end

    subgraph "Event"
        SubscriptionEvent{Event} -->|uses| EventBus((Event Bus))
        RateEvent{Event} -->|uses| EventBus((Event Bus))
    end

    Client[Client] -->|interacts| App(((APP)))

    classDef main fill:#FF2800;
    classDef ultra fill:#FC8000;
    classDef orange fill:#FC8000;
    classDef pink fill:#FF9CD6;
    classDef blue fill:#0000FF;
    classDef yellow fill:#E8BA36;
    classDef func fill:#4EADE5;
    classDef carrot fill:#E35C5C;
    classDef face fill:#700FC1;
    classDef cherry fill:#960039;
    classDef green fill:#39FF14 
    classDef rust fill:#B04721; 
    classDef blues fill:#663DFE; 
    classDef pale fill:#3B898A; 
    classDef pack fill:#C418A2; 
    
    class Web pack;
    class SubscriberEntity,ExchangeRate face;
    class SubscriptionService,RateService cherry;
    class ExchangeRateProvider,Repository,EmailSender,SubscriptionServiceInterface,RateServiceInterface,RateServiceInterface,ExchangeRateServiceInterface func;
    class FileStore pale;
    class EventBus,RateEvent,SubscriptionEvent green;
    class App pack;
    class MessageEntity,TopicEntity blues;
    class SubscriptionHandlers,RateHandlers rust;
    class Logger blue;
    class Env yellow;
    class Client pink; 
   
```