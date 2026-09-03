# broken-hexagon

> An **intentionally broken** hexagonal architecture — a demonstration of
> [`hexago validate`](https://github.com/padiazg/hexago).
>
> Part of the **Gophercon LATAM 2026** presentation.

[HexaGo](https://github.com/padiazg/hexago) scaffolds Go projects with hexagonal
architecture (Ports & Adapters). This repository is a small, compiling Go
service that violates the architecture **on purpose**, so you can see what
`hexago validate` catches.

## What's broken

Everything compiles, tests pass, and `go run main.go run` works. The bugs are
architectural, not syntactic — exactly the kind of rot that accumulates when
nobody enforces the dependency rule:

| # | Where | Violation | Caught by `validate`? |
|---|-------|-----------|:---------------------:|
| 1 | `internal/core/domain/product.go:28` | The domain entity calls `redis.InvalidateProduct` directly — the domain reaches into a secondary adapter | ✅ |
| 2 | `internal/core/services/products/products.go` | The service depends on the concrete `*database.ProductRepository` instead of the `ProductRepository` port in `internal/core/ports` | ✅ |
| 3 | `internal/adapters/secondary/database/product_repository.go` | A secondary adapter imports `primary/http/dto` (a presentation type) — outbound adapter leaks into inbound | ⚠️ [known gap](https://github.com/padiazg/hexago/issues/59) |
| 4 | same file | The repository talks to the `redis` adapter directly, with no port in between | ⚠️ [known gap](https://github.com/padiazg/hexago/issues/58) |

The `⚠️` rows are planted on purpose: adapter-to-adapter imports are currently
not flagged because [`validateAdapterDependencies` is a no-op](https://github.com/padiazg/hexago/issues/58)
— a discussion point about what the validator does (and doesn't) check today.

## Triggering the validation

From the project root:

```console
$ hexago validate
🔍 Validating project: broken-hexagon
   Module: github.com/padiazg/broken-hexagon
   Adapter style: primary-secondary
   Core logic: services

📋 Validation Results:
✓ Domain directory exists
✓ Core logic directory exists
✓ Inbound adapters directory exists
✓ Outbound adapters directory exists
✓ Config directory exists
✓ Adapters follow dependency rules
✓ Using primary for inbound adapters
✓ Using secondary for outbound adapters
✓ Using services for business logic

✗ Domain imports external package: github.com/padiazg/broken-hexagon/internal/adapters/secondary/redis in internal/core/domain/product.go
✗ Services imports adapter: github.com/padiazg/broken-hexagon/internal/adapters/secondary/database in internal/core/services/products/products.go

📊 Summary:
   ✓ Passed: 9
   ⚠️  Warnings: 0
   ✗ Errors: 2

❌ Validation FAILED
Error: validation failed with 2 error(s)
```

It also works from anywhere with an explicit working directory:

```console
$ hexago validate --working-directory /path/to/broken-hexagon
```

### CI gate

`hexago validate` exits **1** on failure, so it drops straight into CI:

```yaml
- name: Hexagonal architecture check
  run: |
    go install github.com/padiazg/hexago/cmd/hexago@latest
    hexago validate
```

## The fix (what valid looks like)

- **Violation 1** — remove the `redis` import from the domain. Cache
  invalidation is a side effect the *application* cares about, so the
  use case should receive an `Invalidator` port, or the secondary adapter
  should handle it internally.
- **Violation 2** — type the service against the port:

  ```go
  type Service struct {
      repo ports.ProductRepository
  }
  ```
- **Violations 3 & 4** — move `ToResponse` out of the repository (the HTTP
  handler already has what it needs), and put a cache port between the
  repository and redis.

Fix the imports and `hexago validate` turns green — the architecture is now
enforced by a command instead of by memory.

## Running the app

```bash
go run main.go run      # long-running service, graceful shutdown on SIGINT
go run main.go version  # build info
make test               # unit tests
```

## Why this exists

Hexagonal architecture is easy to adopt and easy to abandon: the moment the
core is in a hurry, it imports the adapter "just this once". `hexago validate`
turns the dependency rule into a check you can run locally and in CI, so the
architecture survives contact with the codebase.

This repo is the negative test case for that talk: a codebase that looks fine
(`go build` is green) but fails the architecture check.

## License

[MIT](LICENSE) © padiazg

## Related

- [hexago](https://github.com/padiazg/hexago) — the scaffolding tool
- [Hexagonal Architecture](https://alistair.cockburn.us/hexagonal-architecture/)
