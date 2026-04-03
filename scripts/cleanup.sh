#!/usr/bin/env bash

set -euo pipefail

if [[ -t 1 ]]; then
	RED='\033[0;31m'
	GREEN='\033[0;32m'
	BLUE='\033[0;34m'
	BOLD='\033[1m'
	RESET='\033[0m'
else
	RED=''
	GREEN=''
	BLUE=''
	BOLD=''
	RESET=''
fi

MINIMAL_MODE=false

print_section() {
	printf '\n%b==> %s%b\n\n' "${BLUE}${BOLD}" "$1" "$RESET"
}

print_success() {
	printf '\n%b%s%b\n' "$GREEN" "$1" "$RESET"
}

print_error() {
	printf '\n%b%s%b\n' "$RED" "$1" "$RESET"
}

render_progress() {
	local completed=$1
	local total=$2
	local width=28
	local percent=0
	local filled=0
	local empty=0
	local filled_bar=''
	local empty_bar=''

	if (( total > 0 )); then
		percent=$((completed * 100 / total))
		filled=$((completed * width / total))
	fi

	empty=$((width - filled))
	filled_bar=$(printf '%*s' "$filled" '')
	empty_bar=$(printf '%*s' "$empty" '')
	filled_bar=${filled_bar// /#}
	empty_bar=${empty_bar// /-}

	printf '\r[%s%s] %3d%%' "$filled_bar" "$empty_bar" "$percent"
}

run_go_tests_minimal() {
	local total_tests
	local completed=0
	local fifo_path
	local line
	local test_status

	total_tests=$(go test -list '^(Test|Example)' ./... 2>/dev/null | grep -E '^(Test|Example)' | wc -l | tr -d '[:space:]')
	total_tests=${total_tests:-0}
	if (( total_tests == 0 )); then
		total_tests=1
	fi

	fifo_path=$(mktemp -u)
	mkfifo "$fifo_path"

	render_progress 0 "$total_tests"

	set +e
	(
		GOTMPDIR="$PWD/.tmp/go-build" go test -json -p 1 ./... > "$fifo_path" 2>&1
	) &
	local test_pid=$!

	while IFS= read -r line; do
		case "$line" in
			*'"Test":'*'"Action":"pass"'*|*'"Test":'*'"Action":"fail"'*|*'"Test":'*'"Action":"skip"'*)
				((completed++))
				if (( completed > total_tests )); then
					completed=$total_tests
				fi
				render_progress "$completed" "$total_tests"
				;;
		esac
	done < "$fifo_path"

	wait "$test_pid"
	test_status=$?
	set -e

	if (( test_status == 0 )); then
		render_progress "$total_tests" "$total_tests"
	fi
	printf '\n'

	rm -f "$fifo_path"

	return "$test_status"
}

usage() {
	printf 'Usage: bash scripts/cleanup.sh [-m]\n'
	printf '  -m    Run tests in minimal mode with a progress bar instead of verbose per-test output.\n'
}

while getopts ':mh' opt; do
	case "$opt" in
		m)
			MINIMAL_MODE=true
			;;
		h)
			usage
			exit 0
			;;
		\?)
			print_error "Unknown option: -$OPTARG"
			usage
			exit 1
			;;
	esac
done

shift $((OPTIND - 1))

if (( $# > 0 )); then
	print_error "Unexpected arguments: $*"
	usage
	exit 1
fi

finish() {
	local exit_code=$?

	if [[ $exit_code -eq 0 ]]; then
		print_success "CABIN-CHAT cleanup completed successfully."
	else
		print_error "CABIN-CHAT cleanup failed. Check the output above for the first error."
	fi
}

trap finish EXIT

mkdir -p "$PWD/.tmp/go-build"

print_section "Cleaning up Cabin Chat"

print_section "Formatting Go project"
go fmt ./...

print_success "Formatting completed."

print_section "Building Go project"
go build ./...

print_success "Build completed."

print_section "Running Go tests"

if [[ "$MINIMAL_MODE" == true ]]; then
	run_go_tests_minimal
else
	GOTMPDIR="$PWD/.tmp/go-build" go test -p 1 -v ./...
fi

print_success "Tests completed successfully."