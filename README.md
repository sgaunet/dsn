[![Go Report Card](https://goreportcard.com/badge/github.com/sgaunet/dsn)](https://goreportcard.com/report/github.com/sgaunet/dsn)
[![Maintainability](https://api.codeclimate.com/v1/badges/390713a97add700e028c/maintainability)](https://codeclimate.com/github/sgaunet/dsn/maintainability)
[![Test Coverage](https://api.codeclimate.com/v1/badges/390713a97add700e028c/test_coverage)](https://codeclimate.com/github/sgaunet/dsn/test_coverage)

# dsn

Tiny library to handle data source name : scheme://user:password@host:port/dbname&sslmode=disable

It's really dumb but useful for me at least.


And now it's a binary that can be used in bash script.

```
$ eval $(dsn setenv --d "pg://login:password@host/mydb?timeout=1000"  --pr DB_)
$ echo $DB_HOST
host
```

# Install

## Option 1

* Download the release
* Install the binary in /usr/local/bin 

## Option 2: With brew

```
brew tap sgaunet/homebrew-tools
brew install sgaunet/tools/dsn
```

## Option 3: Docker image

Possibility to copy the binary by using the docker image

```
FROM sgaunet/dsn:latest as builder

FROM ....
COPY --from builder /dsn /usr/bin/dsn
```

# Development


This project is using :

* golang 1.19+
* [task for development](https://taskfile.dev/#/)
* docker
* [docker buildx](https://github.com/docker/buildx)
* docker manifest
* [goreleaser](https://goreleaser.com/)

The docker image is only created to simplify the copy of dsn in another docker image.


