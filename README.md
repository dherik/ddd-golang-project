# DDD Go Project

## Introduction

Welcome to DDD Go Project, a dynamic and innovative project built with Go programming language. This initiative is driven by a passion for DDD (Domain Driven Design), Go and test automation.

## How to run

### With Docker Compose

Just run:

```sh
docker compose up
```

To build the application and down the containers:

```sh
docker compose build
docker compose down
```

### Without Docker Compose

Run the databaase:

```sh
docker compose up postgres
```

And run the application using the **Makefile**:

```sh
make build run
```

### How to run the tests

To run the integration tests:

```sh
make integration-test
```

## All itens considered for the project

Below is the extensive list of items I am considering in the project and that I classified as essential for a good Go project. The idea is that each project decision is well-founded. There are still many things to be adjusted, so consider that this project is always evolving.

- Using Docker Compose to facilitate development.
- README.md with all the necessary information to run the project.
- Unit tests.
- Integration tests using the testify and Docker Test libraries to control scenarios.
- REST API endpoints following market standards (using plural nouns, avoiding verbs, employing hierarchy, adopting correct standards for representing dates, etc.).
- Secure endpoints with authentication (JWT).
- Using [conventional commit messages](https://www.conventionalcommits.org/en/v1.0.0/).
- Use of structured logging (Go slog).
- Use of Makefile
- Versioned files with the collection of endpoints, ready to import on [Bruno](https://www.usebruno.com)
- `launch.json` file for Visual Studio Code already configured for debugging the application
- Short names for Go packages.
- Use of dependency injection pattern.

What is coming next? See the [project backlog](https://github.com/users/dherik/projects/1/views/1?layout=board).

## Structure

```
.
├── docker-compose.yml      # use for development
├── Dockerfile
├── docs                    # general documentation
│   └── bruno               # Bruno collections ready to use
│       └── DDD Example
│           ├── bruno.json
│           ├── Create task.bru
│           ├── environments
│           │   └── Local.bru
│           ├── Find task.bru
│           ├── Find task by id.bru
│           └── Login.bru
├── init_ddl.sql            # DDL (SQL to create tables, FKs, etc) for database
├── init_dml.sql            # DML (SQL data) for database
├── internal
│   ├── app                 # application layer (DDD)
│   │   ├── api             # api code
│   │   │   ├── login.go
│   │   │   ├── routes.go   # all routing (endpoints) code 
│   │   │   ├── service.go
│   │   │   └── task.go
│   │   └── server.go
│   ├── domain              # domain layer (DDD)
│   │   ├── task.go
│   │   ├── user.go
│   │   └── user_test.go
│   └── infrastructure      # infrastructure layer (DDD)
│       └── persistence     # persistence code
│           ├── postgresql.go
│           ├── task_memory,go      # memory implementation of task repository
│           ├── task_postgresql.go  # postgre implementation of task repository
│           └── user.go
├── main
├── main.go
├── Makefile
├── README.md
└── tests
    └── integration             # integration tests code 
        ├── setup               # code used to setup the integration tests
        │   ├── database.go
        │   ├── login.go
        │   └── server.go
        └── task_test.go        # integration tests for task

```
<!--
## Introduction

1. `cmd/`: This directory contains the application's entry point (main.go), where you configure and start your application.

1. `internal/app/`: Defines the jobs the software is supposed to do and directs the expressive domain objects to work out problems. The tasks this layer is responsible for are meaningful to the business or necessary for interaction with the application layers of other systems. **This layer is kept thin**. It does not contain business rules or knowledge, but only coordinates tasks and delegates work to collaborations of domain objects in the next layer down. It does not have state reflecting the business situation, but it can have state that reflects the progress of a task for the user or the program.

1. `internal/domain/`: Responsible for representing concepts of the business, information about the business situation, and business rules. State that reflects the business situation is controlled and used here, even though the technical details of storing it are delegated to the infrastructure. This layer is the heart of business software.

1. `internal/infrastructure/`: Defines the jobs the software is supposed to do and directs the expressive domain objects to work out problems. The tasks this layer is responsible for are meaningful to the business or necessary for interaction with the application layers of other systems. **This layer is kept thin**. It does not contain business rules or knowledge, but only coordinates tasks and delegates work to collaborations of domain objects in the next layer down. It does not have state reflecting the business situation, but it can have state that reflects the progress of a task for the user or the program.
    - `persistence/` contains the concrete repository implementations (e.g., UserRepository and TaskRepository) that interact with the database.
    - `messaging/` includes integration code for RabbitMQ.
    - `external/` contains code for integrating with external services like the weather API.
1. `internal/config/`: This directory manages application configuration, including loading configuration settings from config.yaml.

With this structure, you maintain a clear separation of concerns between domain logic, application services, and infrastructure code. You can apply DDD principles to your domain entities, services, and repositories while keeping your codebase organized and easy to understand.

Dependencies in a DDD Service: the Application layer depends on Domain and Infrastructure, and Infrastructure depends on Domain, but Domain doesn't depend on any layer:

```
┌─────────────────────┐
│ Application layer   ├───────────┐
└─────────┬───────────┘           │
          │                       │
          │           ┌───────────▼─────────┐
          │           │ Domain model layer  │
          │           └───────────▲─────────┘
          │                       │
          │                       │
┌─────────▼────────────┐          │
│ Infrastructure layer ├──────────┘
└──────────────────────┘
```

![Detailed organization](DDD-Layers-implemented.png)

### About `domain/repository.go`

The repository.go file typically defines repository interfaces for interacting with your domain entities, such as the User and Task entities. These interfaces provide a contract that concrete repository implementations must adhere to, allowing you to abstract the data access layer and facilitate testing and flexibility.

Each of the interfaces should be implemented by concrete repository implementations, which are responsible for interacting with the database or data store. These implementations will handle the specific database queries and operations required to fulfill the contract defined by the interfaces.

By defining repository interfaces and using them in your application, you decouple the business logic from the data access layer, making your code more modular and testable. Additionally, it allows you to switch between different data storage solutions (e.g., PostgreSQL, MySQL, NoSQL databases) with minimal code changes by implementing the same interfaces for each storage backend. 

### Domain vs Application: business logic

The business logic is primarily associated with the domain layer, but it can also exist in the application layer. The exact division between the two layers can vary depending on the architectural approach you choose and your application's specific requirements. Here's a more detailed explanation:

Domain Layer:
- The domain layer is the primary home for the core business logic of your application. It contains the domain entities (e.g., User and Task) and their associated behaviors and invariants.
- Business rules that are tightly related to the structure and behavior of domain entities should reside in this layer. For example, validation rules for ensuring that a task's due date is in the future or that a user's email address is unique within the system.
- The domain layer is focused on modeling the problem domain and maintaining the integrity of domain entities.

Application Layer:
- The application layer serves as an orchestrator of the domain logic. It contains use cases or services that coordinate the interactions between domain entities and enforce high-level business rules.
- Complex business workflows that involve multiple domain entities often reside in the application layer. For instance, the process of creating a task, checking user permissions, and notifying users about task updates can be orchestrated in the application layer.
- The application layer can also include cross-cutting concerns like logging, error handling, and transaction management.

## Reference

- https://learn.microsoft.com/en-us/dotnet/architecture/microservices/microservice-ddd-cqrs-patterns/ddd-oriented-microservice
- https://www.baeldung.com/hexagonal-architecture-ddd-spring
- https://github.com/ddd-crew/free-ddd-learning-resources
- https://threedots.tech/post/common-anti-patterns-in-go-web-applications/
- https://threedots.tech/post/things-to-know-about-dry/
- https://github.com/ThreeDotsLabs/wild-workouts-go-ddd-example

-->