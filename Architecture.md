# Go Backend — Production Grade Architecture & Folder Structures
### (Raw `net/http`, industry-standard patterns)

---

## STEP 1: Why Folder Structure Matters in Go

Go doesn't force any folder structure (unlike Rails/Django which are opinionated). This means:
- Bad structure → tightly coupled, hard-to-test, hard-to-scale codebase.
- Industry has converged on a few **proven patterns** based on project size & team structure.

The right structure depends on:
1. **App size** (small tool vs large product)
2. **Team size** (1 dev vs 50 devs across teams)
3. **Growth expectation** (will it split into microservices later?)

👉 Let's start with the most basic/common baseline structure used almost everywhere in Go, then build up to bigger patterns.

---

## STEP 2: The Base Convention — `cmd/` + `internal/` + `pkg/`

Almost every serious Go project (regardless of architecture style) uses this root convention, based on the community-adopted **golang-standards/project-layout**:

```
myapp/
├── cmd/                    → entrypoints (main.go files)
│   └── api/
│       └── main.go
├── internal/                → private application code (can't be imported by other projects)
│   └── ...
├── pkg/                     → public/shared code (can be imported by other projects)
│   └── ...
├── api/                     → API contracts (OpenAPI/Swagger specs, proto files)
├── configs/                 → config files (yaml/env templates)
├── deployments/             → Docker, k8s manifests
├── scripts/                 → build/deploy/migration scripts
├── test/                    → extra external test data/helpers
├── go.mod
└── go.sum
```

**Why `internal/` matters:** Go's compiler enforces that anything inside `internal/` can ONLY be imported by code within the same module — this is Go's built-in way of hiding implementation details. It's the single most important folder in production Go apps.

👉 Now let's see how `internal/` is organized differently depending on the architecture style. First up: the simplest — **Layered Monolith**.

---

## STEP 3: Architecture 1 — Layered Monolith (Simple/Classic)

Best for: small-to-medium apps, small teams, early-stage products, single deployable binary.

This organizes code **by technical layer** (handler → service → repository), similar to MVC.

```
myapp/
├── cmd/
│   └── api/
│       └── main.go              → wires everything together, starts http.Server
├── internal/
│   ├── handler/                 → HTTP handlers (parse request, call service, write response)
│   │   ├── user_handler.go
│   │   └── order_handler.go
│   ├── service/                 → business logic
│   │   ├── user_service.go
│   │   └── order_service.go
│   ├── repository/              → DB access layer (queries)
│   │   ├── user_repository.go
│   │   └── order_repository.go
│   ├── model/                   → structs/entities shared across layers
│   │   ├── user.go
│   │   └── order.go
│   ├── middleware/               → auth, logging, recovery, CORS middlewares
│   │   ├── auth.go
│   │   └── logger.go
│   ├── router/                   → route registration (maps URL → handler)
│   │   └── router.go
│   └── config/                   → app config loading (env vars, yaml)
│       └── config.go
├── pkg/                          → reusable utils (hashing, jwt, validators)
│   ├── jwtutil/
│   └── validator/
├── migrations/                   → DB migration files
├── go.mod
└── go.sum
```

**Flow:** `Router → Handler → Service → Repository → Database`

**Pros:** Simple, easy to onboard new devs, fast to build.
**Cons:** As the app grows, `service/` and `repository/` folders become huge dumping grounds — everything is coupled by layer, not by feature, making it hard to find "everything related to Orders."

👉 This limitation is exactly why **Modular Monolith** exists. Let's look at that next.

---

## STEP 4: Architecture 2 — Modular Monolith (Feature/Domain based)

Best for: medium-to-large apps, growing teams, apps that *might* split into microservices later.

Instead of organizing by technical layer, you organize **by business domain/feature (module)** — each module has its own handler+service+repo internally.

```
myapp/
├── cmd/
│   └── api/
│       └── main.go
├── internal/
│   ├── user/                     → USER MODULE (self-contained)
│   │   ├── handler.go
│   │   ├── service.go
│   │   ├── repository.go
│   │   ├── model.go
│   │   └── routes.go
│   ├── order/                    → ORDER MODULE (self-contained)
│   │   ├── handler.go
│   │   ├── service.go
│   │   ├── repository.go
│   │   ├── model.go
│   │   └── routes.go
│   ├── payment/                  → PAYMENT MODULE
│   │   ├── handler.go
│   │   ├── service.go
│   │   ├── repository.go
│   │   ├── model.go
│   │   └── routes.go
│   ├── shared/                    → shared kernel — used by multiple modules
│   │   ├── database/
│   │   ├── middleware/
│   │   └── errors/
│   └── config/
├── pkg/
├── migrations/
├── go.mod
```

**Key rule:** Modules should NOT directly import each other's internal details — they should talk via well-defined interfaces or an internal event system. This keeps modules loosely coupled (mimicking how microservices would talk, but still in one deployable binary).

**Pros:**
- Easy to find everything about one feature ("Order module") in one place.
- If needed, a module can later be **extracted into its own microservice** with minimal refactor (huge advantage).
- Good middle ground between simplicity and scalability.

**Cons:** Slightly more setup overhead than plain layered monolith; requires discipline to avoid modules reaching into each other's internals.

👉 This module-per-domain idea is actually a mini preview of microservices. Let's now go fully distributed.

---

## STEP 5: Architecture 3 — Microservices

Best for: large-scale systems, multiple independent teams, need independent scaling/deployment of different parts.

Each service is its **own separate Go module/repo** (or in a monorepo but independently deployable), with its own DB, own `main.go`, own lifecycle.

### Option A: Multi-repo (each service = separate repository)

```
user-service/
├── cmd/api/main.go
├── internal/
│   ├── handler/
│   ├── service/
│   ├── repository/
│   └── model/
├── pkg/
├── go.mod

order-service/
├── cmd/api/main.go
├── internal/...
├── go.mod

payment-service/
├── cmd/api/main.go
├── internal/...
├── go.mod
```
Each service has an identical internal structure to a small monolith — because from its own point of view, it IS a small monolith (this is a key insight: microservices = many small monoliths talking to each other).

### Option B: Monorepo (all services in one repo, but independently deployable)

```
company-backend/
├── services/
│   ├── user-service/
│   │   ├── cmd/main.go
│   │   ├── internal/
│   │   └── go.mod              → separate go.mod per service
│   ├── order-service/
│   │   ├── cmd/main.go
│   │   ├── internal/
│   │   └── go.mod
│   └── payment-service/
│       ├── cmd/main.go
│       ├── internal/
│       └── go.mod
├── pkg/                          → SHARED across all services (published as internal Go module)
│   ├── logger/
│   ├── jwtutil/
│   ├── proto/                    → shared gRPC/protobuf contracts
│   └── errors/
├── deployments/
│   ├── docker-compose.yml
│   └── k8s/
│       ├── user-service/
│       ├── order-service/
│       └── payment-service/
└── go.work                       → Go workspace file (links multiple go.mod together for local dev)
```

**Note:** `go.work` (Go Workspaces, since Go 1.18) is the modern way to develop multiple modules together locally in a monorepo without needing `replace` directives in go.mod.

**Cross-service communication typically uses:**
- **gRPC** (most common in Go microservices — fast, strongly typed via protobuf)
- **REST/HTTP** (simpler, more universal)
- **Message Queues** (Kafka, RabbitMQ, NATS) for async/event-driven communication

**Pros:** Independent scaling, independent deployment, teams can own services fully, fault isolation.
**Cons:** Operational complexity (service discovery, distributed tracing, network failures, data consistency across services), needs strong DevOps/infra investment.

👉 Now let's look at an internal architecture pattern used WITHIN any of the above (monolith or each microservice) to keep business logic clean and testable — **Clean/Hexagonal Architecture**.

---

## STEP 6: Architecture Pattern 4 — Clean Architecture / Hexagonal (Ports & Adapters)

This isn't an alternative to monolith/microservices — it's a way to structure the **internal** folder of ANY of the above for maximum testability and independence from frameworks/DB.

Core idea: **Business logic (domain) should not depend on frameworks, DB drivers, or HTTP — dependencies point INWARD toward the domain.**

```
internal/
├── domain/                      → CORE — entities + business rules (pure Go, no external deps)
│   ├── user.go                  → User struct + business rules/validation
│   └── user_repository.go       → INTERFACE (port) — defines what repo must do, no implementation
├── usecase/                     → application logic — orchestrates domain rules (also called "service" or "interactor")
│   └── user_usecase.go
├── adapter/                     → ADAPTERS — implementations of the interfaces defined in domain/
│   ├── http/                    → HTTP handler adapter (raw net/http handlers)
│   │   └── user_handler.go
│   ├── repository/               → DB adapter — implements domain's repository interface
│   │   └── postgres_user_repository.go
│   └── external/                 → 3rd party API adapters (payment gateway, email service)
│       └── email_adapter.go
└── config/
```

**Dependency direction:** `adapter → usecase → domain` (never the reverse). Domain knows NOTHING about HTTP or Postgres — you could swap Postgres for MongoDB or REST for gRPC without touching domain/usecase code.

**Pros:** Maximum testability (domain logic tested with zero mocks needed for DB/HTTP), framework-agnostic, very maintainable long-term.
**Cons:** More boilerplate/interfaces upfront — overkill for small apps or prototypes.

👉 This pattern is often combined with Modular Monolith or Microservices (each module/service internally follows Clean Architecture). Let's summarize when to use what.

---

## STEP 7: Which Architecture to Choose? (Decision Guide)

| Situation | Recommended Architecture |
|---|---|
| Small app, MVP, solo dev, prototype | **Layered Monolith** |
| Growing product, single team, might scale later | **Modular Monolith** |
| Multiple teams, need independent scaling/deploys | **Microservices** |
| Need long-term maintainability & heavy business logic | **Clean/Hexagonal Architecture** (inside monolith or each microservice) |
| Large company, many domains, high traffic | **Modular Monolith → gradually extracted into Microservices** (most common real-world migration path) |

**Real industry pattern:** Most successful companies (Shopify, Uber early days, etc.) **start with a Modular Monolith**, and only split modules into microservices **once a specific module has a genuine scaling/team-ownership need** — not upfront. This avoids the "premature microservices" trap that adds huge operational overhead too early.

---

## STEP 8: Common Supporting Folders (used across ALL architectures)

Regardless of which architecture style you pick, these are near-universal in production Go backends:

| Folder | Purpose |
|---|---|
| `migrations/` | SQL migration files (using tools like `golang-migrate`, `goose`, `atlas`) |
| `configs/` or `internal/config/` | Environment-based config loading (dev/staging/prod) |
| `deployments/` | Dockerfile, docker-compose.yml, k8s manifests |
| `scripts/` | Build, deploy, seed-data, codegen scripts |
| `docs/` or `api/` | OpenAPI/Swagger specs, architecture decision records (ADRs) |
| `.github/workflows/` | CI/CD pipelines (lint, test, build, deploy) |
| `Makefile` | Common dev commands (`make run`, `make test`, `make migrate`) |

---

## STEP 9: Raw `net/http` Specific Notes (no framework)

Since you're using raw `net/http` (not Gin/Echo/Fiber), a few structure-specific notes:

- **Routing:** Since Go 1.22, `net/http`'s built-in `ServeMux` supports method-based routing and path parameters natively (`GET /users/{id}`) — many production apps now skip third-party routers entirely.
- **Middleware pattern:** Implemented as function wrappers around `http.Handler` — chained together in `router/` or `middleware/` folder (e.g., `Logger(Auth(CORS(handler)))`).
- **Server setup:** `main.go` in `cmd/api/` typically does: load config → connect DB → build dependencies (repo→service→handler) → register routes → start `http.Server` with proper timeouts (`ReadTimeout`, `WriteTimeout`, `IdleTimeout`) and graceful shutdown (`context` + `signal.NotifyContext`).
- **Dependency Injection:** Go has no built-in DI framework — production apps do **manual constructor injection** (pass dependencies as function/struct parameters) in `main.go`. This is intentional in idiomatic Go — keeps things explicit and simple over "magic" DI containers.

---

## STEP 10: Final Summary

| Architecture | Organized By | Best For | Complexity |
|---|---|---|---|
| **Layered Monolith** | Technical layer (handler/service/repo) | Small apps, MVPs | Low |
| **Modular Monolith** | Business domain/feature | Growing apps, single team | Medium |
| **Microservices** | Independent services | Large scale, multiple teams | High |
| **Clean/Hexagonal** | Dependency direction (domain-centric) | Long-term maintainability, complex business logic | Medium-High (as an overlay on any of the above) |

**One-line rule of thumb:**
> Start with Modular Monolith + Clean internal structure. Extract to Microservices only when a specific module truly needs independent scaling or team ownership.