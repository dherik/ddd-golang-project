# DDD Go Project

A small REST API for managing tasks, built with Go and organized around Domain-Driven Design (DDD) layers. The project focuses on two things above everything else: **DDD design** and **integration testing** with Testcontainers. It doubles as a template for new Go services that want a pragmatic starting point: explicit layers, dependency injection, authentication, persistence, messaging, and automated tests.

> **Using this repository as a template?** Start with the [Using this project as a template](#using-this-project-as-a-template) checklist before writing your own domain code.

## Introduction

The project is a simple REST API (with authentication) for managing tasks. Its focus is the **DDD design** — explicit layers, the dependency rule, and a domain model free of HTTP and database concerns — and the **integration testing setup**, which exercises the full API against real PostgreSQL and RabbitMQ instances provisioned with Testcontainers.

The DDD strategies applied are intentionally basic, so consider this a simplified implementation of the design. More advanced concepts can be layered on top.

The code still has small issues and things to improve. You can follow the [project backlog](https://github.com/users/dherik/projects/1) to see what is coming next.

## Using this project as a template

When starting a new project from this repository, make it yours before adding features. The most reusable parts are the layer structure (see [Architecture](#architecture)) and the integration test harness (see [Tests](#tests)):

1. **Rename the module**: update `go.mod` and imports, for example `go mod edit -module github.com/your-user/your-project`, then replace `github.com/dherik/ddd-golang-project` across the code.
2. **Replace the domain**: delete the task/user examples once you have your own entities, value objects, and repository interfaces in `internal/domain`.
3. **Replace the database schema**: rewrite `init_ddl.sql` and `init_dml.sql` for your tables. Keep the same layering, and do not expose database rows outside `internal/infrastructure/persistence`.
4. **Replace the credentials**: change the PostgreSQL credentials in `docker-compose.yml`, `config.yaml`, and `init_ddl.sql`.
5. **Set your own JWT secret**: use the `JWT_SECRET` environment variable instead of the development secret committed in `config.yaml`.
6. **Review the example endpoints**: keep, adapt, or remove `/tasks`, `/users`, and the RabbitMQ listener, depending on what your service needs.
7. **Remove what you do not use**: the Bruno collection, k6 load tests, Ollama/Continue scripts, and the RabbitMQ integration can all be deleted without affecting the core API.
8. **Update this README**: remove template-specific links and document your own domain, endpoints, and configuration.
9. **Add what a template cannot decide for you**: CI, license, deployment manifests, and environment-specific configuration.

## Architecture

The project follows the classic DDD dependency rule: inner layers define the contracts, outer layers implement them. The domain layer has no knowledge of HTTP, SQL, or messaging. This layer separation is one of the two main focuses of the project: the domain stays independent while HTTP, persistence, and messaging live in replaceable outer layers.

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

- **Domain** (`internal/domain`): entities, validation rules, repository interfaces, and the password-hashing and clock ports. No third-party dependencies.
- **Application** (`internal/app`): HTTP handlers, services/use cases, DTOs, routes, and message consumers. Depends on domain and (through interfaces) infrastructure.
- **Infrastructure** (`internal/infrastructure`): repository implementations, RabbitMQ client, password hashing, and external integrations. Depends on domain.
- **Config** (`internal/config`): loads `config.yaml` and applies environment variable overrides.

### Where to add new code

| What you are adding | Where it goes |
|---|---|
| Entity, value object, domain rule | `internal/domain` |
| Repository interface | `internal/domain` (next to the entity) |
| Use case / service | `internal/app/api` (or a dedicated use case package as the project grows) |
| HTTP handler, request/response DTO, route | `internal/app/api` |
| Repository implementation | `internal/infrastructure/persistence` |
| Password hashing or clock adapter | `internal/infrastructure/security` (port in `internal/domain`) |
| External API client | `internal/infrastructure/external` |
| Message consumer | `internal/app/messaging` + `internal/infrastructure/messaging` |

### Project structure

```sh
.
├── config.yaml             # Application configuration
├── docker-compose.yml      # Local development environment
├── Dockerfile              # Multi-stage, non-root production image
├── init_ddl.sql            # Schema (tables, constraints, indexes)
├── init_dml.sql            # Seed data
├── docs
│   └── bruno               # Bruno collection ready to import
├── internal
│   ├── app
│   │   ├── api             # Handlers, services, routes, auth helper
│   │   ├── messaging       # Message listener and handlers
│   │   └── server.go       # Application wiring and server startup
│   ├── config              # Configuration loader
│   ├── domain              # Domain entities, ports, and repository interfaces
│   └── infrastructure
│       ├── external        # External integrations
│       ├── messaging       # RabbitMQ client
│       ├── persistence     # PostgreSQL repositories
│       └── security        # bcrypt password hasher
├── tests
│   ├── integration         # Testcontainers-based integration tests
│   └── load                # k6 load tests
├── main.go
├── Makefile
└── README.md
```

## Getting started

### Prerequisites

- [Docker](https://www.docker.com) (required for Docker Compose and integration tests)
- [Go](https://go.dev) 1.27.2 or newer. [GVM](https://github.com/moovweb/gvm) is a practical way to manage Go versions.
- Optional: [k6](https://github.com/grafana/k6) for load tests and [Bruno](https://www.usebruno.com) for the API collection.

### Run everything with Docker Compose

```sh
docker compose build
docker compose up
```

To stop the containers:

```sh
docker compose down
```

### Run the application locally with Docker services

Start PostgreSQL and RabbitMQ:

```sh
docker compose up postgres rabbitmq
```

Then build and run the application with the Makefile:

```sh
make build run
```

This method makes debugging easier because the Go process runs on the host.

### Configuration

Configuration lives in `config.yaml` and can be overridden with environment variables. This keeps the file safe to commit while allowing deployment-specific values.

| Config key | Environment variable(s) | Default |
|---|---|---|
| `http.port` | `HTTP_PORT` | `3333` |
| `database.host` | `DB_HOST`, `DATABASE_HOST` | `localhost` |
| `database.port` | `DB_PORT`, `DATABASE_PORT` | `5432` |
| `database.user` | `DB_USER`, `DATABASE_USER` | `pguser` |
| `database.password` | `DB_PASSWORD`, `DATABASE_PASSWORD` | `pgpassword` |
| `database.name` | `DB_NAME`, `DATABASE_NAME` | `dddtasks` |
| `jwt.secret` | `JWT_SECRET` | development value in `config.yaml` |
| `broker.host` | `BROKER_HOST`, `RABBITMQ_HOST` | `localhost` |
| `broker.port` | `BROKER_PORT`, `RABBITMQ_PORT` | `5672` |
| `broker.user` | `BROKER_USER`, `RABBITMQ_USER` | `guest` |
| `broker.password` | `BROKER_PASSWORD`, `RABBITMQ_PASSWORD` | `guest` |

Always set `JWT_SECRET` in production. Never commit a real secret.

### Make targets

| Target | Description |
|---|---|
| `make build` | Builds the `main` binary |
| `make run` | Builds and runs the application against local services |
| `make unit-test` | Runs unit tests only (`-short`, integration tests are skipped) |
| `make integration-test` | Runs all tests with coverage, including Testcontainers-based integration tests |
| `make test` | Runs every test verbosely |
| `make vulnerability-check` | Runs `govulncheck` against the module |
| `make upgrade-deps-latest-patch` | Updates dependencies to the latest patch versions |
| `make upgrade-deps-latest-version` | Updates dependencies to their latest versions |
| `make load-test-create-tasks` | Runs the k6 task creation load test |
| `make clean` | Removes the built binary |

## Tests

Integration testing is the other main focus of the project. The test suite runs against the real HTTP server and real infrastructure, so it validates the wiring between the DDD layers instead of mocking it away. Unit tests cover the domain and application rules in isolation.

### Unit tests

```sh
make unit-test
```

or:

```sh
go test -short ./...
```

`-short` is important: integration tests check `testing.Short()` and skip themselves, so a green short run does not execute them.

### Integration tests

The integration tests start PostgreSQL and RabbitMQ with [Testcontainers](https://golang.testcontainers.org), start the real HTTP server, and exercise the API end to end. Docker must be running.

```sh
make integration-test
```

or the integration suite only:

```sh
go test -count=1 -v ./tests/integration/...
```

Both setup and teardown are embedded in the test code, so the suite can also be executed from an IDE.

### Race detector

```sh
go test -race -short ./...
```

## API

The API is served at `http://localhost:3333` by default.

### Authentication

1. `POST /login` with a username and password.
2. Receive a JWT signed with `HS256`, valid for 72 hours.
3. Send it on protected endpoints as `Authorization: Bearer <token>`.

The token contains `userId`, `name`, `email`, `sub`, `iat`, and `exp`. The seeded `admin` user (`some_password`) exists only for local development.

```sh
TOKEN=$(curl -s -X POST http://localhost:3333/login \
  -H 'Content-Type: application/json' \
  -d '{"username":"admin","password":"some_password"}' | jq -r .token)
```

### Endpoints

| Method | Path | Authentication | Description |
|---|---|---|---|
| `POST` | `/login` | Public | Exchanges credentials for a JWT. |
| `POST` | `/users` | JWT | Creates a user. Returns `409` for duplicate username or email. |
| `GET` | `/tasks` | JWT | Lists the authenticated user's tasks, newest first, one page at a time. Supports optional, inclusive `startDate` and `endDate` query parameters in RFC 3339 format, plus `limit` (default `20`, maximum `100`) and `offset` (default `0`). |
| `GET` | `/tasks/:id` | JWT | Returns one task owned by the authenticated user, or `404`. |
| `POST` | `/tasks` | JWT | Creates a task for the authenticated user. The owner always comes from the token; the body only needs a `description`. |
| `GET` | `/api/health` | Public | Liveness endpoint. |

Task ownership is enforced in the persistence layer: a user can never read or create tasks for another user. Asking for someone else's task returns `404`, so task IDs are not leaked.

### Pagination

`GET /tasks` returns one page at a time:

```json
{
  "items": [
    {
      "id": 1,
      "description": "Complete project proposal",
      "userId": "1",
      "createdAt": "2024-02-15T10:59:01.054Z"
    }
  ],
  "hasMore": false
}
```

- `items` is always an array, ordered newest first by `createdAt` and then by descending `id`.
- `hasMore` is `true` only when at least one task exists beyond the returned page.
- `limit` defaults to `20` and accepts `1` to `100`. A value outside that range returns `400`; it is never clamped.
- `offset` defaults to `0` and accepts any non-negative integer. An offset past the end returns an empty page with `hasMore: false`.
- `startDate` and `endDate` are optional, independent and inclusive; omit either to leave that side of the range unbounded.
- An empty parameter value is treated as absent.
- There is no total count and no page number.

Pagination is offset-based. A task created between page requests shifts later rows down, so the last task of one page can reappear at the top of the next. To freeze a walk, send the `createdAt` of the first page's newest task as `endDate` on every subsequent request.

### Examples

```sh
# Create a task for the authenticated user
curl -s -X POST http://localhost:3333/tasks \
  -H "Authorization: Bearer $TOKEN" \
  -H 'Content-Type: application/json' \
  -d '{"description":"Write the README"}'

# Find one page of tasks in a date range (RFC 3339)
curl -s "http://localhost:3333/tasks?startDate=2024-02-10T00:00:00Z&endDate=2024-02-25T23:59:59Z&limit=20&offset=0" \
  -H "Authorization: Bearer $TOKEN"

# Find the newest tasks without a date filter
curl -s "http://localhost:3333/tasks?limit=10" \
  -H "Authorization: Bearer $TOKEN"

# Find a task by ID
curl -s http://localhost:3333/tasks/1 \
  -H "Authorization: Bearer $TOKEN"
```

### Error responses

Errors use a consistent JSON body:

```json
{
  "timestamp": "2026-10-17T12:00:00Z",
  "error": {
    "code": 404,
    "message": "task not found"
  }
}
```

| Status | Meaning |
|---|---|
| `200` | Request succeeded |
| `201` | Resource created |
| `400` | Invalid payload, query parameter, or path parameter |
| `401` | Missing, invalid, or expired token, or invalid credentials |
| `404` | Resource not found (or not owned by the authenticated user) |
| `409` | Duplicate username or email |
| `500` | Unexpected server error |

## Load tests

Install [k6](https://github.com/grafana/k6), start the stack with `docker compose up`, and run:

```sh
make load-test-create-tasks
```

The test ramps up to 50 virtual users, logs in as the seeded admin user, and creates tasks. It finishes after 3 minutes.

## Security notes

- Set `JWT_SECRET` through the environment; the committed value is for local development only.
- Passwords are hashed with bcrypt at the library default cost.
- Tasks are scoped to the authenticated user in every task endpoint.
- Duplicate usernames and emails are rejected on the database level, not only in application code.
- `make vulnerability-check` runs `govulncheck` against the module and its dependencies.

## Technical features

The DDD design and the integration testing setup are the two areas this project focuses on. The remaining features exist to support them.

### Architecture and code quality

- Domain models independent from the other layers: persistence rows and transport DTOs do not leak into the domain.
- Classic DDD layering with the dependency rule pointing inward (Application → Domain, Infrastructure → Domain, Domain → nothing).
- Dependency injection through constructors and small interfaces.
- Short, focused names for Go packages.
- User-owned resources enforced in the persistence layer using the user ID from the JWT.

### Testing

- Unit tests using [Testify](https://github.com/stretchr/testify).
- Integration tests using Testify suites and [Testcontainers](https://golang.testcontainers.org) to provision PostgreSQL and RabbitMQ.
- All integration test setup is embedded in the code, so the suite runs directly from your IDE.
- `-short` separates fast unit tests from the full integration suite.

### API and security

- REST API following common conventions: plural nouns, no verbs, hierarchy, and RFC 3339 dates.
- Endpoints protected with JWT authentication; the token carries `userId`, `name`, and `email`.
- Consistent JSON error responses with meaningful HTTP status codes.
- Passwords hashed with bcrypt in an infrastructure adapter at the library default cost.
- `govulncheck` integrated in the Makefile for dependency vulnerability scanning.

### Developer experience and tooling

- Docker Compose to facilitate local development.
- Simple steps to have the service up and running on a developer machine.
- Multi-stage, non-root production Docker image.
- Configuration through `config.yaml` with environment variable overrides.
- Well-organized `Makefile` covering build, run, tests, dependency upgrades (latest and latest patch), vulnerability check, and load tests.
- Organized `README.md` with everything necessary to run and test the project.
- Versioned [Bruno collection](docs/bruno/) ready to import.
- Visual Studio Code `launch.json` already configured for debugging.
- Conventional commit messages and a pull request template.
- Structured logging with Go `slog`, request IDs, and latency.

## Known gaps

These are intentional next steps rather than surprises:

- No refresh tokens or logout endpoint.
- No rate limiting, body size limit, or configurable HTTP server timeouts yet.
- Integration tests currently use a fixed port (`3333`) and stop only the database container on teardown.

See the [project backlog](https://github.com/users/dherik/projects/1/views/1?layout=board) for the current plan.


## Technical references

- [Microsoft - Microservices architecture, DDD, and CQRS](https://learn.microsoft.com/en-us/dotnet/architecture/microservices/microservice-ddd-cqrs-patterns/ddd-oriented-microservice)
- [Baeldung - Hexagonal Architecture, DDD, and Spring](https://www.baeldung.com/hexagonal-architecture-ddd-spring)
- [Free DDD Learning Resources](https://github.com/ddd-crew/free-ddd-learning-resources)
- [Common Anti-patterns in Go Web Applications](https://threedots.tech/post/common-anti-patterns-in-go-web-applications/)
- [Things to Know About DRY (Don't Repeat Yourself)](https://threedots.tech/post/things-to-know-about-dry/)
- [Wild Workouts Go DDD Example by ThreeDotsLabs](https://github.com/ThreeDotsLabs/wild-workouts-go-ddd-example)
