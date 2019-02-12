# GitHub Follower API

## Overview

GitHub Follower API utilizes the `followers` endpoint of the GitHub v3 REST API in order to return a collection of usernames that follow a given user, as well as the followers of their followers.

**Important:** The JSON data returned will only contain the usernames of the top 100 followers up to 4 layers deep.

Read more about the `follower` endpoint of the Github v3 REST API here: <https://developer.github.com/v3/users/followers/>

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

`-t` for test dependencies.

### Configuration

See the `/config/config.go` constants.

There you can configure:

```bash
DefaultPort
MaxFollowerCount
MaxFollowerDepth
```

### Building & Running

```bash
go build
./github-follower-api
```

or

```bash
./github-follower-api 3000
```

to run on a different port.

The `DefaultPort` is 8080.

## Tests

### Running Tests

In order to run the tests recursively for this project:

```bash
go test ./...
```

Tests are organized into subtests. You can pass -v to `go test` to see their labels and table-driven test details.

### Structure

All unit tests are located next to their respective production code. I.e., `api_test.go` is next to `api.go`.
