# envcheck

`envcheck` is an infrastructure-focused Go CLI for CI/CD preflight validation.

It validates required environment variables and checks external dependencies (HTTP endpoints and TCP services) before deploy steps run.

## Why this project

- Reproducible automation-first CLI design.
- CI-friendly exit codes (`0` success, `1` failed checks, `2` usage/runtime errors).
- Single release binary artifact for Linux amd64.
- Build-time version metadata injected with `-ldflags`.

## Project structure

```text
cmd/envcheck/main.go      # CLI flags, output formatting, exit codes
internal/checks/env.go    # required environment validation
internal/checks/http.go   # HTTP health checks with timeout
internal/checks/tcp.go    # TCP connectivity checks with timeout
internal/version/version.go
```

## Usage

Run locally:

```bash
go test ./...
go build -o envcheck ./cmd/envcheck
./envcheck --required-env DB_HOST,DB_PORT --check-http https://example.com/health --check-tcp 127.0.0.1:5432
```

JSON output for machine consumers:

```bash
./envcheck --required-env DB_HOST --json
```

Show build version:

```bash
./envcheck --version
```

Repeatable checks and custom timeout:

```bash
./envcheck \
  --check-http https://example.com/health \
  --check-http https://example.com/ready \
  --check-tcp 127.0.0.1:5432 \
  --timeout 3s
```

## CI workflow

**File:** `.github/workflows/ci.yml`  
**Trigger:** push to `main`

CI executes:

1. `go build ./...`
2. `go test ./...`

## CD workflow

**File:** `.github/workflows/cd.yml`  
**Trigger:** GitHub Release created

CD executes:

1. `golangci-lint`
2. `go test ./...`
3. Build `envcheck-linux-amd64` from `./cmd/envcheck`
4. Inject version, commit, and build timestamp via `-ldflags`
5. Upload the binary to the release with `softprops/action-gh-release`

## Example CI preflight command

```bash
./envcheck \
  --required-env DB_HOST,DB_PORT,REDIS_ADDR \
  --check-http https://internal-api.example.com/health \
  --check-tcp redis.example.com:6379 \
  --timeout 2s
```
