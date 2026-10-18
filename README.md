# phys-demo

Demo artifacts for the Phys ecosystem. One directory per benchmark
problem; each demo keeps its own provenance, derivation notes, and module
arrangement. See `../phys-system/AGENTS.md` for the collective
architecture and operating rules.

## Layout

Every Phys derivation project is an independent Go module under the
canonical repository prefix `github.com/PithomLabs/phys-demo` and carries
its own `go.work`, so the standard command `go test ./...` works from
the project root in any context. There is no shared workspace file; a
project's `go.work` is the only workspace context it uses.

- `problem_8a/inkling/` — Go module
  `github.com/PithomLabs/phys-demo/problem_8a/inkling`; its `go.work`
  lists the sibling modules it imports (`phys-artifact`, `phys-math`,
  `phys-lib`), so those resolve from the local multi-repository checkout.
- `problem_8a/muse/` — Go module
  `github.com/PithomLabs/phys-demo/problem_8a/muse`; self-contained
  (`go.work` uses only `.`).
- `problem_8a/muse+deepseek/` — Go module
  `github.com/PithomLabs/phys-demo/problem_8a/muse-deepseek`
  (the directory name `muse+deepseek` is kept; `+` is not a legal
  character in a Go module path, hence the sanitized `muse-deepseek`
  module segment). Self-contained (`go.work` uses only `.`).

## Testing a single project

Run from the project root, with no environment overrides:

    cd problem_8a/muse
    go test ./...

This genuinely discovers and tests every package in that project's
module. `go vet ./...` and `go build ./...` work the same way.

## Testing all demo projects

Run from this directory:

    ./scripts/test-all.sh

The script resolves the repository root itself, so it also works when
invoked by absolute path from any other working directory.

The runner runs the standard `go test -count=1 ./...` from each project
root. Because every project owns its context through its own `go.work`,
the runner needs no workspace-mode special casing. It prints the project
being tested, continues with the remaining projects after a failure, and
exits nonzero if any suite fails. A listed directory without `go.mod` is
reported as an error rather than skipped silently.

`go test ./...` at this repository root is **not** the mechanism for
testing the projects: this directory is not itself a Go module, and Go's
`./...` pattern does not cross nested module boundaries in any case.
The runner is the supported cross-project entry point.

### Adding another demo project

1. Ensure the directory contains a `go.mod` with a unique module path
   under `github.com/PithomLabs/phys-demo`.
2. Add a `go.work` (`go` directive plus `use .`), listing every sibling
   module the project imports alongside `.`.
3. Add one entry to the `projects` array in `scripts/test-all.sh`.

### Caveats

- `problem_8a/inkling` resolves `phys-artifact`, `phys-math`, and
  `phys-lib` from the sibling checkouts in `../../..` via its `go.work`.
  Cloning `phys-demo` without those sibling repositories leaves inkling
  unresolvable by design (fail-fast with an explicit error); `muse` and
  `muse+deepseek` remain fully self-contained.
- Do not run `go mod tidy` in the projects until the sibling modules
  publish a consistent version set; tidy may try to rewrite `go.mod`
  from versions that only exist in local checkouts.
