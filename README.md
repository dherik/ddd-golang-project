# DDD Go Project

## Introduction

Welcome to the DDD Go Project, a dynamic project built with the Go programming language. This initiative is driven by a passion for DDD (Domain-Driven Design), Go, and test automation. For now the project is just a simple REST API (with authentication) to manage tasks.

The code still has a lot of small issues and new things to implement. You can check the [Backlog of the project](https://github.com/users/dherik/projects/1) to follow what's comming up.

## How to run

You can run the project using:
- Docker Compose;
- Using the `make` command and initializing just the database (also with Docker Compose). This method is easier to debug the application.

Both methods are explained below.

### With just Docker Compose

Just run:

```sh
docker compose up
```

To build the application or bring down the containers:

```sh
docker compose build
docker compose down
```

### With Docker Compose + Makefile

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

Below is the extensive list of items I am considering in the project and that I classified as essential for a good Go project. The idea is that each project will have each decision well-founded. As mentioned before, there are still many things to be adjusted, so consider that this project is always evolving.

- Using Docker Compose to facilitate development
- Organized README.md with all the necessary information to run the project
- Unit tests
- Integration tests using the testify and Docker Test libraries to control test scenarios
- REST API endpoints following market standards (using plural nouns, avoiding verbs, employing hierarchy, adopting correct standards for representing dates, etc.)
- Secure endpoints with authentication (JWT)
- Using [conventional commit messages](https://www.conventionalcommits.org/en/v1.0.0/)
- Use of structured logging (Go slog)
- Use of Makefile
- Versioned files with the collection of endpoints, ready to import on [Bruno](https://www.usebruno.com)
- `launch.json` file for Visual Studio Code already configured for debugging the application
- Short names for Go packages
- Use of dependency injection pattern
- Domain models totally independent from other layers, not using domain objects on persistence (database) or application (API) layers.
- Use of `gosec` to find security vulnerabilities

What is coming next? See the [project backlog](https://github.com/users/dherik/projects/1/views/1?layout=board).

## Structure

```
.
├── docker-compose.yml      # Used for development
├── Dockerfile
├── docs                    # General documentation
│   └── bruno               # Bruno collections ready to import and use
│       └── DDD Example
│           ├── bruno.json
│           ├── Create task.bru
│           ├── environments
│           │   └── Local.bru
│           ├── Find task.bru
│           ├── Find task by id.bru
│           └── Login.bru
├── init_ddl.sql            # DDL (SQL to create tables, FKs, etc) for the database
├── init_dml.sql            # DML (SQL data) for database
├── internal
│   ├── app                 # Application layer (DDD)
│   │   ├── api             # API code
│   │   │   ├── login.go
│   │   │   ├── routes.go   # All routing (endpoints) code 
│   │   │   ├── service.go
│   │   │   └── task.go
│   │   └── server.go
│   ├── domain              # Domain layer (DDD)
│   │   ├── task.go
│   │   ├── user.go
│   │   └── user_test.go
│   └── infrastructure      # Infrastructure layer (DDD)
│       └── persistence     # Persistence code
│           ├── postgresql.go
│           ├── task.go     # Implementation of the task repository
│           └── user.go     # Implementation of the user repository
├── main
├── main.go
├── Makefile
├── README.md
└── tests
    └── integration             # Integration tests code 
        ├── setup               # Code used to setup the integration tests
        │   ├── database.go     # Code for database initilization for integration tests
        │   ├── login.go        
        │   ├── setup.go        # Common code for integration tests setup
        │   └── server.go       # Code for HTTP server initialization for integration tests
        ├── task_test.go        # Integration tests for task
        └── user_test.go        # Integration tests for user

```

With this structure, you maintain a clear separation of concerns between domain logic, application services, and infrastructure code. You can apply DDD principles to your domain entities, services, and repositories while keeping your codebase organized and easy to understand.

Dependencies in a DDD Service: the Application layer depends on Domain and Infrastructure, and Infrastructure depends on Domain, but Domain doesn't depend on any layer:

```
┌──────────────────────┐
│ Application layer    ├──────────┐
└─────────┬────────────┘          │
          │           ┌───────────▼────────┐
          │           │ Domain model layer │
          │           └───────────▲────────┘
┌─────────▼────────────┐          │
│ Infrastructure layer ├──────────┘
└──────────────────────┘
```

## Technical references

- [Microsoft - Microservices architecture, DDD, and CQRS](https://learn.microsoft.com/en-us/dotnet/architecture/microservices/microservice-ddd-cqrs-patterns/ddd-oriented-microservice)
- [Baeldung - Hexagonal Architecture, DDD, and Spring](https://www.baeldung.com/hexagonal-architecture-ddd-spring)
- [Free DDD Learning Resources](https://github.com/ddd-crew/free-ddd-learning-resources)
- [Common Anti-patterns in Go Web Applications](https://threedots.tech/post/common-anti-patterns-in-go-web-applications/)
- [Things to Know About DRY (Don't Repeat Yourself)](https://threedots.tech/post/things-to-know-about-dry/)
- [Wild Workouts Go DDD Example by ThreeDotsLabs](https://github.com/ThreeDotsLabs/wild-workouts-go-ddd-example)
