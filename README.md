# atlacp

[![Build](https://github.com/gemyago/atlacp/actions/workflows/build-flow.yml/badge.svg)](https://github.com/gemyago/atlacp/actions/workflows/build-flow.yml)
[![Coverage](https://raw.githubusercontent.com/gemyago/atlacp/test-artifacts/coverage/golang-coverage.svg)](https://htmlpreview.github.io/?https://raw.githubusercontent.com/gemyago/atlacp/test-artifacts/coverage/golang-coverage.html)

An MCP (Model Context Protocol) interface for Atlassian products (Jira and Bitbucket).

## Project structure

* [cmd/server](./cmd/server) is a main entrypoint to start the server
* [cmd/jobs](./cmd/jobs) is a main entrypoint to start jobs
* [internal/api/http](./internal/api/http) - includes http routes related stuff
  * [internal/api/http/v1routes.yaml](./internal/api/http/v1routes.yaml) - OpenAPI spec for the api routes. HTTP layer is generated with [apigen](github.com/gemyago/apigen)
* `internal/app` - place to add application layer code (e.g business logic).
* `internal/services` - lower level components are supposed to be here (e.g database access layer e.t.c).

## Project Setup

Please have the following tools installed: 
* [direnv](https://github.com/direnv/direnv) 
* [gobrew](https://github.com/kevincobain2000/gobrew#install-or-update)

Install/Update dependencies: 
```sh
# Install
go mod download
go install tool

# Update:
go get -u ./... && go mod tidy
```

### Build dependencies

This step is required if you plan to work on the build tooling. In this case please make sure to install:
* [pyenv](https://github.com/pyenv/pyenv?tab=readme-ov-file#installation).

```sh
# Install required python version
pyenv install -s

# Setup python environment
python -m venv .venv

# Reload env
direnv reload

# Install python dependencies
pip install -r requirements.txt
```

If updating python dependencies, please lock them:
```sh
pip freeze > requirements.txt
```

## Development

### Lint and Tests

Run all lint and tests:
```bash
make lint
make test
```

Run specific tests:
```bash
# Run once
go test -v ./internal/api/http/v1controllers/ --run TestHealthCheck

# Run same test multiple times
# This is useful to catch flaky tests
go test -v -count=5 ./internal/api/http/v1controllers/ --run TestHealthCheck

# Run and watch. Useful when iterating on tests
gow test -v ./internal/api/http/v1controllers/ --run TestHealthCheck
```
### Run local API server:

```bash
# Regular mode
go run ./cmd/server/ start

# Watch mode (double ^C to stop)
gow run ./cmd/server/ start
```