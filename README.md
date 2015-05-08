# Boilerplate

> The way we [Go](http://golang.org/)

Boilerplate provides a simple tool which you can use to provision a
new Go project with a Makefile, Dockerfile, and associated files.

Boilerplate revolves around 3 concepts, each of which are used to set up your new project:

* `repository`: the name of the source control repository _(e.g. github.com)_
* `namespace`: the name of the organization/group in the repository _(e.g. zulily)_
* `project`: the name of the binary _(e.g. fizzbuzz)_

Boilerplate makes/enforces several assumptions about the structure and conventions of a Go project.  Among them:

* The project contains a single binary (aka `package main`).  The name of this binary is the same as that of the project.
* All dependencies are managed via [godep](https://github.com/tools/godep).
* All builds are compiled in a Docker container, using a pinned version of Go (v1.4.2 at the time of writing)
* Binaries are compiled as *true* static binaries, with no cgo or dynamically-linked networking packages.
* Binaries have the `main.BuildSHA` var set to the latest `git` SHA in the repo.  This is accomplished using the Go linker's `-ldflags -X` option.
* A Docker image is created for the resulting binary, which `exec`s the binary as the entrypoint.
* The Docker image uses the naming convention `<namespace>/<project>`, and is tagged with the latest `git` SHA in the repo.

## Prerequisites

`boilerplate` depends on having the following installed on the host machine:

* `git`
* `make`
* `go`
* `Docker` (or `boot2docker` for OSX users)

## Quick Start

`boilerplate` may be invoked with no arguments, in which case you will be
interactively prompted for the `repository`, `namespace`, and `project` names.

    $ go get github.com/core/boilerplate
    $ boilerplate

    Enter the name of git repository (e.g. github.com): github.com
    Enter the namespace in the repository (e.g. zulily): dcarney
    Enter the name of the project (e.g. fizzbuzz): whizbang

    GOPATH is: /home/dcarney/go
    Creating a new project at: /home/dcarney/go/src/github.com/dcarney/whizbang
    Creating new: .dockerignore
    Creating new: .gitignore
    Creating new: Dockerfile
    Creating new: Makefile
    Creating new: main.go
    Initializing git repo
    Initializing godeps
    Done

## Example

Values for the `repository`, `namespace`, and `project` can also be supplied using command line flags:

    $ boilerplate -repository=github.com -namespace=zulily project=fizzbuzz

    GOPATH is: /home/dcarney/go
    Creating a new project at: /home/dcarney/go/src/foobar/zulily/fizzbuzz
    Creating new: .dockerignore
    Creating new: .gitignore
    Creating new: Dockerfile
    Creating new: Makefile
    Creating new: main.go
    Initializing git repo
    Initializing godeps
    Done

The resulting project can be compiled, linted, and "Dockerized" using the supplied `Makefile` targets:

    $ cd $GOPATH/src/foobar/zulily/fizzbuzz
    $ make
    building binary for fizzbuzz...

    $ make lint
    linting code...
    main.go:8:6: exported type Foobar should have comment or be unexported
    main.go:15:2: can probably use "var slice []string" instead

    $ make dockerize
    building binary for fizzbuzz...
    running tests for fizzbuzz...
    building Docker image 'zulily/fizzbuzz'...

    $ docker images | grep fizzbuzz
    zulily/fizzbuzz                 HEAD                  4f679023d74c        5 seconds ago      1.94 MB

See the generated `Makefile` in your `boilerplate`-created project for more details and build targets

## Limitations and known issues

Boilerplate was developed on, and is used daily on, various Linux OSes.  It has been known to work on OSX with `boot2docker`. It has not been used on Windows. (pull requests welcome!)
