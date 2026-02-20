# go-cicd-demo

A minimal Go CLI that prints a greeting. Used to demonstrate CI/CD with GitHub Actions.

## The Go Project

- **main.go** — Defines `Greet(name string)` and a `main` that prints the greeting (or "Hello, World!" if no argument is given).
- **main_test.go** — Unit tests for `Greet`.

Run locally:

```bash
go build .
./go-cicd-demo
./go-cicd-demo Alice
go test ./...
```

## CI Workflow

**File:** `.github/workflows/ci.yml`

**Trigger:** Every push to the `main` branch.

**Steps:**

1. Check out the repository.
2. Set up Go 1.25 (using `actions/setup-go@v5`).
3. Run `go build ./...`.
4. Run `go test ./...`.

If any step fails, the workflow fails and the push is not considered verified.

## CD Workflow

**File:** `.github/workflows/cd.yml`

**Trigger:** When a GitHub Release is **created** (e.g. via the repo’s Releases page or `gh release create`).

**Steps:**

1. Check out the repository.
2. Set up Go 1.25.
3. Lint the code with **golangci-lint** (`golangci/golangci-lint-action@v6`).
4. Run `go test ./...`.
5. Build a Linux amd64 binary: `go build -o go-cicd-demo .` (with `GOOS=linux`, `GOARCH=amd64`).
6. Upload the binary to the newly created release using **softprops/action-gh-release@v2**.

The workflow has `contents: write` so `GITHUB_TOKEN` can attach the binary to the release.

**Creating a release:** In the repo, go to **Releases** → **Create a new release**, choose a tag (e.g. `v1.0.0`), publish the release. The CD workflow runs and attaches the `go-cicd-demo` binary to that release.

## GitHub Actions Used

| Action | Purpose |
|--------|---------|
| [actions/checkout@v4](https://github.com/actions/checkout) | Check out the repository. |
| [actions/setup-go@v5](https://github.com/actions/setup-go) | Install and cache a Go toolchain. |
| [golangci/golangci-lint-action@v6](https://github.com/golangci/golangci-lint-action) | Run golangci-lint (CD only). |
| [softprops/action-gh-release@v2](https://github.com/softprops/action-gh-release) | Upload the built binary to the GitHub Release (CD only). |
