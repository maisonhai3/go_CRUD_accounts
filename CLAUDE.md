# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## What this is

A learning/practice repo for Go backend fluency, centered on an account-CRUD HTTP service backed by SQLite. `README.md` frames the goals as a self-assessment in Vietnamese — the user is drilling specific patterns (handler signatures, context propagation, trust boundaries, graceful shutdown). When helping, favor explaining *why* a pattern is used over just making code compile; the point is internalization, not throughput.

## Commands

```bash
go run .                      # run the account service on :8080 (root package main)
go build -o accountCRUD.exe . # build the service binary
go test ./...                 # run tests (none exist yet — README lists tests as a requirement)
go test -run TestCreateAccount ./handlers   # run a single test by name once tests exist
go vet ./...                  # static checks
```

The standalone learning sketches each have their own `main` and run independently:

```bash
go run ./leak                 # goroutine-leak / context-cancellation experiments
```

## Build caveat

`example.go` at the repo root is `package main` alongside `main.go` but is an **incomplete scratch file with deliberate syntax errors** (e.g. a bare `var`, `fmt.print()` with no import). It breaks `go build .` / `go run .` for the root package. Treat it as a scratchpad — if a build of the main service is needed, that file must be fixed or excluded first. Confirm intent before "fixing" it; it may be a work-in-progress exercise.

## Architecture

Three-layer dependency flow, one direction only: `main.go` → `handlers/` → `repositories/`. The transport layer depends on the data layer; the data layer knows nothing about HTTP.

- **`main.go`** — wires `repositories.DBHandler` into `handlers.Handler`, registers routes on `http.ServeMux` using Go 1.22+ method+path patterns (`"POST /accounts"`, `"GET /accounts/{id}"`), configures the `http.Server` with explicit timeouts, and runs a **graceful shutdown** loop on SIGINT/SIGTERM via `srv.Shutdown(ctx)`.
- **`handlers/`** — HTTP transport. `Handler.Repo` is the only dependency. Each handler derives a per-request `context.WithTimeout` from `r.Context()`.
- **`repositories/`** — `DBHandler` wraps `*sql.DB`. All queries take `ctx` and use `QueryRowContext`/`QueryContext`/`ExecContext`.

### Trust-boundary DTO pattern (the core idea of this repo)

Data is never decoded straight into the domain struct. Each layer has its own type and mapping is explicit:

`CreateAccountRequest` (handler input DTO) → validated → `CreateAccountParams` (repo input) → `Account` (domain) → `AccountResponse` (handler output DTO).

Mapping helpers live in `handlers/accountHandlers.go` (`toAccountResponse`, `toAccountResponses`). When adding a field, expect to touch all the relevant DTOs deliberately — that duplication is the design, not an accident. README §"Can it be attacked?" explains the rationale (input type and domain type are different concerns).

### Conventions that must be preserved

- **Money is `int64` minor units, never `float64`** (see schema comment in `repositories/init_db.go`). Do not introduce floating-point money.
- **Timestamps are stored as RFC3339 `TEXT`** in SQLite and parsed back to `time.Time` in `scanAccount`. Reads scan into `string` then `time.Parse(time.RFC3339, ...)`.
- **Soft delete**: every read filters `deleted_at IS NULL`; deletion sets the column rather than removing the row.
- **Currency is whitelisted** in `handlers/utils.go` (`isValidCurrency`, currently `VND`/`USD`). Extend there, not inline.
- **Input hardening on writes** (`CreateAccount`): `http.MaxBytesReader` caps body size, `decoder.DisallowUnknownFields()` rejects extra fields, server forces `balance = 0` regardless of client input.
- **`scanAccount`** takes a `rowScanner` interface so the same scan logic serves both `QueryRowContext` (single row) and `QueryContext` rows (list).
- **Not-found signaling**: repo translates `sql.ErrNoRows` into the package sentinel `repositories.ErrAccNotFound`; handlers match it with `errors.Is` to return 404.

### SQLite driver

Uses `modernc.org/sqlite` — a **pure-Go** driver (no CGO, no C toolchain needed). The driver name passed to `sql.Open` is `"sqlite"` (not `"sqlite3"`), and it is registered via the blank import `_ "modernc.org/sqlite"` in `repositories/init_db.go`. The DSN enables WAL mode and a busy timeout. The DB file `accounts.db` is created in the working directory on first run.

## Learning sketches (not part of the service)

- **`leak/`** — `package main` exploring goroutine leaks, `context` cancellation/timeout, `sync.WaitGroup` + buffered channels + `select` fan-in.
- **`streaming_exp/`** — `package streaming_exp`, streaming-parse of a large JSON array from an HTTP body via `json.Decoder.Token()` + `More()` + per-element `Decode()` (constant memory, no full-array load).
