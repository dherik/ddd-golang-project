# DDD Go Project

## Introduction

Welcome to the DDD Go Project, a dynamic project built with the Go programming language. This initiative is driven by a passion for DDD (Domain-Driven Design), Go, and test automation. 

For now, the project is just a simple REST API (with authentication) for managing tasks. The DDD strategies applied are still very basic, so consider this project as a simplified implementation of this design. However, I plan to introduce more advanced concepts in the future.

The code still has a lot of small issues and new things to implement. You can check the [Backlog of the project](https://github.com/users/dherik/projects/1) to follow what's comming up.

## How to run

Prerequisites:
- [Docker](https://www.docker.com)
- [Go](https://go.dev) 1.27.2 or newer. I recommend to use [GVM](https://github.com/moovweb/gvm) to manage the Go versions

You can run the project using choosing on the following methods:
- Docker Compose;
- Using the `make` command and initializing just the database (also with Docker Compose). This method is easier to debug the application.

Both methods are explained below.

### With just Docker Compose

Just run:

```sh
docker compose build
docker compose up
```

To bring down the containers:

```sh
docker compose down
```

### With Docker Compose + Makefile

Run the database and the RabbitMQ broker:

```sh
docker compose up postgres rabbitmq
```

And run the application using the **Makefile**:

```sh
make build run
```

### How to run the tests

To run the unit tests:

```sh
make unit-test
```

To run the integration tests (requires Docker):

```sh
make integration-test
```

## Run K6 load test

Install the latest version of [K6](https://github.com/grafana/k6/) to be able to execute the tests.

After install K6, just run `docker compose up` in one terminal and `make load-test-create-tasks` on another terminal. The test will finish in 3 minutes. You can use the same process for the other `load-test-*`.

## Technical features

Below is the extensive list of items that are implemented in the project that I classified as essential for a good Go project and make the developer's life easier as possible. As mentioned before, there are still many things to be adjusted, so consider that this project is always evolving.

- Using Docker Compose to facilitate development
- Simple steps to have the service up and running in the developer machine
- Organized `README.md` file with all the necessary information to run and test the project
- Unit tests
- Integration tests using the testify to control test scenarios and TestContainers to set up the infrastructure (Postgres and RabbitMQ)
- You can execute integration tests directly from your favorite IDE, as all necessary setup configurations are embedded within the code.
- REST API endpoints following best practice standards (using plural nouns, avoiding verbs, employing hierarchy, adopting correct standards for representing dates, etc.)
- Secure endpoints with authentication (JWT)
- Using [conventional commit messages](https://www.conventionalcommits.org/en/v1.0.0/)
- Use of structured logging using Go `slog`
- Very well organized `Makefile`
- Versioned files with the [collection of endpoints](docs/bruno/), ready to import on [Bruno](https://www.usebruno.com)
- `launch.json` file for Visual Studio Code already configured for debugging the application
- Use of short names for Go packages
- Use of dependency injection pattern
- Domain models totally independent from other layers, not using domain objects on persistence (database) or application (API) layers.
- Use of `govulncheck` to find security vulnerabilities
- Use of Go upgrade tool to upgrade the dependencies to the latest version
- Pull Request template to provide a baseline standard of informational quality and organizational rigor
- Use of Local LLM ([LM Studio](https://lmstudio.ai) + [Llama 3](https://huggingface.co/meta-llama/Meta-Llama-3-8B) model) to help in the initial project setup
- Use of [Continue](https://www.continue.dev/) extension as AI code assistant, using the models [Llama 3 8B](https://ollama.com/library/llama3) for chat and [Starcoder 2 3B](https://ollama.com/library/starcoder2:3b) for code generation, using [Ollama](https://ollama.com) as model server.

What is coming next? See the [project backlog](https://github.com/users/dherik/projects/1/views/1?layout=board).

## Structure

```sh
.
├── config.yaml             # Application configuration
├── docker-compose.yml      # Used for development
├── Dockerfile
├── docs                    # General documentation
│   └── bruno               # Bruno collections ready to import and use
│       ├── bruno.json
│       ├── Create task.bru
│       ├── Create user.bru
│       ├── Find task.bru
│       ├── Find task by id.bru
│       ├── Login.bru
│       └── environments
│           └── Local.bru
├── init_ddl.sql            # DDL (SQL to create tables, FKs, etc)
├── init_dml.sql            # DML (seed data) for database
├── internal
│   ├── app                 # Application layer (DDD)
│   │   ├── api             # API handlers, services, and routes
│   │   │   ├── login.go
│   │   │   ├── routes.go
│   │   │   ├── service.go
│   │   │   ├── task.go
│   │   │   └── user.go
│   │   ├── messaging       # Message listeners and handlers
│   │   │   ├── handlers.go
│   │   │   └── listener.go
│   │   └── server.go
│   ├── config              # Configuration loader
│   │   └── config.go
│   ├── domain              # Domain layer (DDD)
│   │   ├── password.go
│   │   ├── password_test.go
│   │   ├── task.go
│   │   ├── task_test.go
│   │   ├── user.go
│   │   └── user_test.go
│   └── infrastructure      # Infrastructure layer (DDD)
│       ├── external
│       │   └── weatherapi.go
│       ├── messaging
│       │   └── rabbitmq
│       │       └── rabbitmq.go
│       └── persistence
│           ├── postgresql.go
│           ├── task.go
│           └── user.go
├── main.go
├── Makefile
├── README.md
├── scripts
│   └── start-ollama-continue.sh
└── tests
    └── integration
        ├── integration_test.go
        ├── setup
        │   ├── broker.go
        │   ├── broker_command.go
        │   ├── database.go
        │   ├── login.go
        │   ├── server.go
        │   └── setup.go
        ├── task_test.go
        └── user_test.go
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

## AI Code assistant

You can setup the AI code assistant installing the [Continue](https://www.continue.dev/) extension, the [Ollama](https://ollama.com) server and following the instructions and executing this [README.md](scripts/README.md) file.

## Technical references

- [Microsoft - Microservices architecture, DDD, and CQRS](https://learn.microsoft.com/en-us/dotnet/architecture/microservices/microservice-ddd-cqrs-patterns/ddd-oriented-microservice)
- [Baeldung - Hexagonal Architecture, DDD, and Spring](https://www.baeldung.com/hexagonal-architecture-ddd-spring)
- [Free DDD Learning Resources](https://github.com/ddd-crew/free-ddd-learning-resources)
- [Common Anti-patterns in Go Web Applications](https://threedots.tech/post/common-anti-patterns-in-go-web-applications/)
- [Things to Know About DRY (Don't Repeat Yourself)](https://threedots.tech/post/things-to-know-about-dry/)
- [Wild Workouts Go DDD Example by ThreeDotsLabs](https://github.com/ThreeDotsLabs/wild-workouts-go-ddd-example)
