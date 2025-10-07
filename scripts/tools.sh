#!/usr/bin/env bash

# shellcheck disable=SC2046

# shellcheck disable=SC2145
echo "tools script launched as: '$@'"

SCRIPT_PATH=$(
    cd "$(dirname "${BASH_SOURCE[0]}")" || exit 1
    pwd
)

UCELL_BANNERATOR_SOURCES_VERSION=$(
    cd "$(dirname "${BASH_SOURCE[0]}")" || exit 1
    v=$(git tag -l --points-at HEAD)
    if [ -n "$v" ]; then
        echo "$v"
    else
        TZ=UTC git show -s --abbrev=12 --date=format-local:%Y%m%d%H%M%S --pretty=format:v0.0.0-%cd-%h
    fi
)

function intro() {
    echo "Script parameters:

* Script was launched from: ${SCRIPT_PATH}
* ucell bannerator service sources version: ${UCELL_BANNERATOR_SOURCES_VERSION}
* Go cache: $(go env GOCACHE)
* HOME: ${HOME} (note when it is empty or set to '/' go will effectively disable build cache!)
* GOPATH: ${GOPATH}
* User ID: ${UID}
* Architecture: $(uname -m)"
}

function generate() {
    run_command go generate ./...
}

function lint() {
    run_command golangci-lint --version
    run_command golangci-lint run --modules-download-mode=vendor
}

function lint_fix() {
    run_command golangci-lint --version
    run_command golangci-lint run --fix --modules-download-mode=vendor
}

function gotest() {
    run_command task build-test-migrator
    run_command task migrate-test-db
    run_command go test -p 1 $(go list ./...) -v -coverprofile coverage.out
    run_command go tool cover -func=coverage.out
}

function gotest_race() {
    run_command task build-test-migrator
    run_command task migrate-test-db
    run_command go test -race $(go list ./...) -v -coverprofile coverage.out
    run_command go tool cover -func=coverage.out
}

function run_command() {
    start=${SECONDS}
    echo "> ${*}"
    (eval "${*}")
    exitcode=$?
    diff=$((${SECONDS} - $start))

    if [ $exitcode -ne 0 ]; then
        echo "> ...execution failed after ${diff}s :("
        exit 1
    fi

    echo "> ...execution succeeded in ${diff}s :)"
}

# Parse CLI parameters.
for i in "$@"; do
    case $i in
    generate)
        intro
        generate
        ;;
    lint)
        intro
        lint
        ;;
    lint-fix)
        intro
        lint_fix
        ;;
    test)
        intro
        gotest
        ;;
    test_race)
        intro
        gotest_race
        ;;
    esac
done
