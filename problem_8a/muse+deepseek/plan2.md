# Plan — Root-Level Test Runner for `phys-demo`

Status: **implemented and locally verified** (root `scripts/test-all.sh` + root `README.md`; see §6). Git ownership/commit pending: `phys-demo/` is not currently a Git repository.

Task: provide one root-level command that tests every intended Go demo module
from `/home/chaschel/Documents/go/phys-demo` in its correct Go module/workspace
context, reports failures accurately, and exits nonzero if any suite fails.
Repository-root `go test ./...` is explicitly **not** the mechanism (Go's
package pattern does not cross nested module boundaries).

## 1. Recon findings (verified, not assumed)

```
phys-demo/                      ← NOT a git repo; no README, go.mod, go.work, scripts/, Makefile
└── problem_8a/
    ├── inkling/                Go module phys-demo/problem_8a/inkling (go 1.23)
    │                           listed in parent workspace → tests pass under ambient GOWORK (2 pkgs ok)
    ├── muse/                   Go sources, but NO go.mod → not a Go module
    └── muse+deepseek/          standalone module github.com/PithomLabs/phys-demo/problem_8a/muse-deepseek (go 1.25.7)
                                NOT in parent workspace; NO local go.work either
                                intended context: GOWORK=off → tests pass (2 pkgs ok + 1 with no test files)
```

- `find` (depth 6) found **exactly two `go.mod`** files and **zero `go.work`,
  `*.sh`, `Makefile`** anywhere under `phys-demo`. Every directory is either a
  module (2) or an ordinary directory (`muse/`).
- Bare root command behavior, both verified:
  - ambient workspace: `pattern ./...: directory prefix . does not contain modules listed in go.work or their selected dependencies`
  - `GOWORK=off`: `pattern ./...: directory prefix . does not contain main module or its selected dependencies`
- `shellcheck`: **not installed** → validation will use `bash -n` only (no new
  dependencies added for linting).
- `AGENTS.md` consulted: workspace root pointer + `phys-system/AGENTS.md`
  (tooling-only change; touches neither `phys`, `phys-gr`, `phys-artifact`,
  `phys-math`, `phys-lib`, nor `phys-trace`).
- Discrepancy flag: the task brief mentions "the existing … local workspace
  configuration" for `muse+deepseek`, but the live tree has **no local
  `go.work`** there. The recorded arrangement is a standalone module tested
  with `GOWORK=off`. It will be kept exactly as is (no local workspace added
  or removed).

## 2. Proposed changes (2 files, tooling only)

### a) New file: `phys-demo/scripts/test-all.sh` (+x)

```bash
#!/usr/bin/env bash
# scripts/test-all.sh — test every supported Go module in phys-demo.
#
# Go's "./..." pattern does not cross nested module boundaries, so a bare
# `go test ./...` from the repository root cannot cover the demo modules.
# This script is the explicit orchestration layer: it runs
# `go test -count=1 ./...` from each module root, in that module's correct
# Go workspace context, and exits nonzero if any suite fails.
#
# Adding a module: add a "<relative-dir>|<mode>" entry to `modules` below.
#   workspace : module is listed in an enclosing go.work (currently the
#               parent workspace /home/chaschel/Documents/go/go.work);
#               tested with ambient workspace discovery.
#   off       : standalone module not listed in any workspace;
#               tested with GOWORK=off (module mode).
# A listed directory without go.mod is reported as a failure,
# never treated as a pass.

set -u -o pipefail

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

modules=(
	"problem_8a/inkling|workspace"
	"problem_8a/muse+deepseek|off"
)

exit_status=0
tested=0
failed=0

for entry in "${modules[@]}"; do
	dir="${entry%%|*}"
	mode="${entry#*|}"

	printf '\n==> %s (%s)\n' "$dir" "$mode"

	if [ ! -f "$repo_root/$dir/go.mod" ]; then
		printf 'ERROR: %s has no go.mod; not a Go module — not tested\n' "$dir" >&2
		failed=$((failed + 1))
		exit_status=1
		continue
	fi

	case "$mode" in
	workspace | off) ;;
	*)
		printf 'ERROR: unknown workspace mode "%s" for %s\n' "$mode" "$dir" >&2
		failed=$((failed + 1))
		exit_status=1
		continue
		;;
	esac

	tested=$((tested + 1))

	run_status=0
	if [ "$mode" = "off" ]; then
		(cd "$repo_root/$dir" && GOWORK=off go test -count=1 ./...) || run_status=$?
	else
		(cd "$repo_root/$dir" && unset GOWORK && go test -count=1 ./...) || run_status=$?
	fi

	if [ "$run_status" -eq 0 ]; then
		printf 'PASS: %s\n' "$dir"
	else
		printf 'FAIL: %s (exit %s)\n' "$dir" "$run_status" >&2
		failed=$((failed + 1))
		exit_status=1
	fi
done

printf '\n==> %d module(s) tested, %d failed\n' "$tested" "$failed"
exit "$exit_status"
```

Design points: repo root resolved from `BASH_SOURCE` (works from any cwd);
explicit module table (documented as the extension point); per-module
workspace context; `unset GOWORK` for the workspace module so a caller's
stray `GOWORK=off` cannot break `inkling`; continue-after-failure; nonzero
exit if any module fails; no new dependencies, no network (all modules and
workspace members are local).

### b) New file: `phys-demo/README.md` (root)

No root README exists, so the brief's "add a section to the root README"
becomes "create a minimal root README containing exactly that section"
(layout + testing), ~40 lines: the root command `./scripts/test-all.sh`, the
per-module context table (`inkling` → enclosing workspace;
`muse+deepseek` → `GOWORK=off`), failure reporting (continues, exits nonzero;
a missing `go.mod` is an error, not a skip), why root `go test ./...` is not
the mechanism (with the two actual error messages above), that `muse/` is not
a module and is not covered, and the 2-step recipe for adding a future module.

## 3. Explicitly NOT changed

`/home/chaschel/Documents/go/go.work`; `problem_8a/inkling/**`;
`problem_8a/muse/**`; `problem_8a/muse+deepseek/**` (including its
`go.mod`/README/plan files); no Makefile (no existing repo convention → avoid
a second implementation); no frozen/shared repos; no commits (not a git repo).

## 4. Validation plan (after implementation)

1. `bash -n scripts/test-all.sh` → syntax check (shellcheck unavailable).
2. From `phys-demo`: `./scripts/test-all.sh` → expect
   `PASS: problem_8a/inkling`, `PASS: problem_8a/muse+deepseek`, summary
   `2 module(s) tested, 0 failed`, exit 0.
3. Invoke once from a different cwd (e.g. `/tmp`) to prove root resolution.
4. Controlled failure: temporarily add
   `problem_8a/inkling/zz_test_all_fail_test.go` with a failing test → run →
   expect `FAIL: …/inkling`, then still `PASS: …/muse+deepseek`
   (continuation), summary `1 failed`, exit 1 → delete the temp file (show
   directory listing before/after so nothing is left behind) → re-run green.
5. Standalone context unchanged: `GOWORK=off go test -count=1 ./...` in
   `muse+deepseek` still green; `go.work` contents re-catted to confirm
   untouched.

## 5. Resolution of open flags (verdict: PASS; implemented)

- Root README: **created new** (none existed) — approved and implemented.
- `muse+deepseek` local-workspace question: resolved by §6 (superseded).
- `muse/` stays untested (not a module); documented as such.

## 6. Decision update (resolved): root runner supersedes the demo-local workspace

- **Superseded:** the earlier option/decision to add a demo-local `go.work`
  inside `problem_8a/muse+deepseek/` so that bare `go test ./...` works from
  that directory.
- **Recorded:** the root-level runner `scripts/test-all.sh` is the supported
  cross-module entry point. `muse+deepseek` remains a standalone module
  tested with `GOWORK=off`; no local `go.work` is added or removed;
  `inkling` keeps its existing parent-workspace context.
- **Consequences:** `go test ./...` inside `problem_8a/muse+deepseek/` still
  requires `GOWORK=off` (as documented in its README); from the `phys-demo/`
  root, use `./scripts/test-all.sh`. No architecture change.
- **Ownership note (pending):** `phys-demo/` is not currently a Git
  repository, so the new runner and README cannot be committed to
  `github.com/PithomLabs/phys-demo` from that directory. Local
  implementation and verification are done; the GitHub deliverable is
  complete only after repository ownership/commit is resolved separately.

### Implemented and verified (this pass)

- Added `phys-demo/scripts/test-all.sh` (executable) and
  `phys-demo/README.md`.
- Normal run from `phys-demo/`: both modules PASS, runner exit 0.
- Invoked from `/tmp/opencode`, and with a caller-exported `GOWORK=off`:
  both modules PASS, exit 0 (root resolution + per-module workspace context
  hold regardless of caller cwd/env).
- Controlled failure (temporary failing test in `problem_8a/inkling`,
  removed afterwards; before/after file-list diff clean):
  `FAIL: problem_8a/inkling (exit 1)`, runner continued to
  `PASS: problem_8a/muse+deepseek`, summary `2 module(s) tested, 1 failed`,
  overall exit 1; re-run after cleanup green (exit 0).
- Standalone demo re-checked from its own module root
  (`GOWORK=off go test -count=1 ./...`): green.
- Parent `/home/chaschel/Documents/go/go.work` re-read: unchanged
  (mtime 2026-10-16 17:32, pre-session). No changes to `muse/`,
  `inkling/**` contents, or any frozen/shared repository.
