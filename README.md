# DDD Go Project

## Introduction

Welcome to DDD Go Project, a dynamic project built with Go programming language. This initiative is driven by a passion for DDD (Domain Driven Design), Go and test automation.

The code still have a lot of small issues and new things to implement. You can check the [Backlog of the project](https://github.com/users/dherik/projects/1) to follow what's comming up.

## How to run

You can run the project using:
- The Docker Compose;
- Using the `make` command and initializing just the database (also with Docker Compose). This method it's easier to debug the application.

Both methods are explained below.

### With Docker Compose

Just run:

```sh
docker compose up
```

To build the application or put down the containers:

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

Below is the extensive list of items I am considering in the project and that I classified as essential for a good Go project. The idea is that each project will have each decision well-founded. As said before, there are still many things to be adjusted, so consider that this project is always evolving.

- Using Docker Compose to facilitate development;
- Organized README.md with all the necessary information to run the project;
- Unit tests;
- Integration tests using the testify and Docker Test libraries to control test scenarios.
- REST API endpoints following market standards (using plural nouns, avoiding verbs, employing hierarchy, adopting correct standards for representing dates, etc.);
- Secure endpoints with authentication (JWT);
- Using [conventional commit messages](https://www.conventionalcommits.org/en/v1.0.0/);
- Use of structured logging (Go slog);
- Use of Makefile;
- Versioned files with the collection of endpoints, ready to import on [Bruno](https://www.usebruno.com)
- `launch.json` file for Visual Studio Code already configured for debugging the application;
- Short names for Go packages;
- Use of dependency injection pattern;
- Domain models totally independent from other layers, not using domain objects on persistence (database) or application (API) layers.

What is coming next? See the [project backlog](https://github.com/users/dherik/projects/1/views/1?layout=board).

## Structure

```
.
├── docker-compose.yml      # use for development
├── Dockerfile
├── docs                    # general documentation
│   └── bruno               # Bruno collections ready to import and use
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
│           ├── task.go     # implementation of task repository
│           └── user.go     # implementation of user repository
├── main
├── main.go
├── Makefile
├── README.md
└── tests
    └── integration             # integration tests code 
        ├── setup               # code used to setup the integration tests
        │   ├── database.go     # code for database initilization for integration tests
        │   ├── login.go        
        │   ├── setup.go        # common code for integration tests setup
        │   └── server.go       # code for HTTP server initialization for integration tests
        ├── task_test.go        # integration tests for task
        └── user_test.go        # integration tests for user

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

- https://learn.microsoft.com/en-us/dotnet/architecture/microservices/microservice-ddd-cqrs-patterns/ddd-oriented-microservice
- https://www.baeldung.com/hexagonal-architecture-ddd-spring
- https://github.com/ddd-crew/free-ddd-learning-resources
- https://threedots.tech/post/common-anti-patterns-in-go-web-applications/
- https://threedots.tech/post/things-to-know-about-dry/
- https://github.com/ThreeDotsLabs/wild-workouts-go-ddd-example