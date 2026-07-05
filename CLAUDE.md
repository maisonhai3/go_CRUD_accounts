# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## What this is

A learning/practice repo for Go backend fluency, centered on an account-CRUD HTTP service backed by SQLite. `README.md` frames the goals as a self-assessment in Vietnamese — the user is drilling specific patterns (handler signatures, context propagation, trust boundaries, graceful shutdown). When helping, favor explaining *why* a pattern is used over just making code compile; the point is internalization, not throughput.

## Commands

The account service is its own Go module rooted at `phase_1_foundations/` — commands run from inside that directory:

```bash
cd phase_1_foundations
go run .                      # run the account service on :8080
go build -o accountCRUD.exe . # build the service binary
go test ./...                 # run tests (none exist yet — README lists tests as a requirement)
go test -run TestCreateAccount ./handlers   # run a single test by name once tests exist
go vet ./...                  # static checks
```

The standalone learning sketches live under `experiments/` (see below) and run independently from the repo root:

```bash
go run ./experiments/leak     # goroutine-leak / context-cancellation experiments
```

## Repo layout

The repo is organized by phase, each its own Go module:

- **`phase_1_foundations/`** — the account CRUD service (`main.go`, `handlers/`, `repositories/`). This is the module described in the Architecture section below.
- **`phase_2_ledger_core/`** — scaffold for the next phase (currently a placeholder `main.go`; not yet implemented).
- **`experiments/`** — learning sketches, grouped here so they don't clutter either phase. `experiments/concurrency_exp/` is its **own Go module** and must be built/vetted from inside that directory; the rest build under whatever module sits above them.

Each phase folder has its own `go.mod`/`go.sum`, so `go build ./...` / `go vet ./...` must be run from inside the relevant phase directory, not the repo root.

## Scratchpad caveat

Some files under `experiments/` are deliberately incomplete exercises. When one won't compile, treat it as a work-in-progress scratchpad and **confirm intent before "fixing" it** — the point of the repo is the user practicing, not a green build. (`example.go`, an earlier root-level scratch of this kind, has since been removed.)

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

Uses `modernc.org/sqlite` — a **pure-Go** driver (no CGO, no C toolchain needed). The driver name passed to `sql.Open` is `"sqlite"` (not `"sqlite3"`), and it is registered via the blank import `_ "modernc.org/sqlite"` in `repositories/init_db.go`. The DSN enables WAL mode and a busy timeout. The DB file `accounts.db` is created in the working directory on first run — i.e. inside `phase_1_foundations/` when run via `cd phase_1_foundations && go run .`.

## Learning sketches (not part of the service)

All under `experiments/`:

- **`experiments/leak/`** — `package main` exploring goroutine leaks, `context` cancellation/timeout, `sync.WaitGroup` + buffered channels + `select` fan-in.
- **`experiments/streaming_exp/`** — `package streaming_exp`, streaming-parse of a large JSON array from an HTTP body via `json.Decoder.Token()` + `More()` + per-element `Decode()` (constant memory, no full-array load). `txn_example.json` is sample input.
- **`experiments/streaming_file/`** — `package streamingfile`, streaming a multipart file upload straight to object storage without buffering the whole file: sniff MIME from a pooled 512-byte head buffer (`sync.Pool`), then `io.MultiReader` the sniffed head back in front of the rest of the body. `aws_stub.go` stands in for a real S3 client so the sketches compile.
- **`experiments/concurrency_exp/`** — **its own Go module** (`golang.org/x/sync`). `package main` at its root covers data races / race conditions; `api_calls/` demoes `errgroup`; `cache_stampede/` demoes `singleflight`. Build and run from inside `experiments/concurrency_exp/`.
