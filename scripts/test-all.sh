#!/usr/bin/env bash
# scripts/test-all.sh — test every derivation project in phys-demo.
#
# Every derivation project is its own Go module with a project-local
# go.work, so the standard command `go test ./...` works from each
# project root in any directory context. This script is only the
# orchestration layer: it runs that standard command from each project
# root, continues with the remaining projects after a failure, and
# exits nonzero if any suite fails.
#
# Adding a project: ensure the directory has a go.mod and a go.work
# (`use .`, plus every sibling module the project imports), then add
# its relative directory to `projects` below.
# A listed directory without go.mod is reported as a failure,
# never treated as a pass.

set -u -o pipefail

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

projects=(
	"problem_8a/inkling"
	"problem_8a/muse"
	"problem_8a/muse+deepseek"
)

exit_status=0
tested=0
failed=0

for dir in "${projects[@]}"; do
	if [ ! -f "$repo_root/$dir/go.mod" ]; then
		printf 'ERROR: %s has no go.mod; not a Go module; not tested\n' "$dir" >&2
		failed=$((failed + 1))
		exit_status=1
		continue
	fi

	tested=$((tested + 1))
	printf '\n==> %s\n' "$dir"

	run_status=0
	(cd "$repo_root/$dir" && unset GOWORK && go test -count=1 ./...) || run_status=$?

	if [ "$run_status" -eq 0 ]; then
		printf 'PASS: %s\n' "$dir"
	else
		printf 'FAIL: %s (exit %s)\n' "$dir" "$run_status" >&2
		failed=$((failed + 1))
		exit_status=1
	fi
done

printf '\n==> %d project(s) tested, %d failed\n' "$tested" "$failed"
exit "$exit_status"
