# GitHub Follower API

## Setup

### Installation

To download this project, either clone the repository or run:

```bash
go get github.com/TGibsonn/github-follower-api
```

### Dependencies

Retrieve dependencies for this project by running:

```bash
go get -d -t ./...
```

### Configuration

TODO.

### Building & Running

```bash
go build
./github-follower-api
```

## Tests

### Running Tests

In order to run the tests recursively for this project:

```bash
go test ./...
```

Tests are organized into subtests. You can pass -v to `go test` to see their labels and table-driven test details.

### Structure

All unit tests are located next to their respective production code. I.e., `api_test.go` is next to `api.go`.
