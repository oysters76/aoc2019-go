#!/bin/sh
# Pre-commit hook for Go puzzle solutions.
#
# For every staged .go file:
#   1. Checks that it's gofmt-formatted (checks the staged version).
#   2. Builds it on its own, so dayN.go and dayN_idomatic.go, which both
#      declare main(), don't clash.
#   3. Runs go vet on it to catch common mistakes.
#
# Skip in an emergency with: git commit --no-verify

files=$(git diff --cached --name-only --diff-filter=ACMR -- '*.go')
[ -z "$files" ] && exit 0

if ! command -v go >/dev/null 2>&1; then
	echo "pre-commit: 'go' not found in PATH" >&2
	exit 1
fi

tmp=$(mktemp -d)
trap 'rm -rf "$tmp"' EXIT
root=$(git rev-parse --show-toplevel)
fail=0

# 1. Formatting
unformatted=""
for f in $files; do
	if [ -n "$(git show ":$f" | gofmt -l)" ]; then
		unformatted="$unformatted $f"
	fi
done
if [ -n "$unformatted" ]; then
	echo "pre-commit: these files are not gofmt-formatted:"
	for f in $unformatted; do echo "  $f"; done
	echo "Fix with:  gofmt -w$unformatted && git add$unformatted"
	fail=1
fi

# 2 and 3. Build and vet each non-test file on its own
for f in $files; do
	case "$f" in *_test.go) continue ;; esac
	dir=$(dirname "$f")
	file=$(basename "$f")

	if ! (cd "$root/$dir" && go build -o "$tmp/out" "$file"); then
		echo "pre-commit: build failed: $f"
		fail=1
		continue
	fi
	if ! (cd "$root/$dir" && go vet "$file"); then
		echo "pre-commit: go vet failed: $f"
		fail=1
	fi
done

if [ "$fail" -ne 0 ]; then
	echo "pre-commit: commit aborted."
	exit 1
fi
exit 0
