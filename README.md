go-sdk
======

![Build Status](https://circleci.com/gh/blend/go-sdk.svg?style=shield&circle-token=:circle-token)

`go-sdk` is our core library of packages. These packages can be composed to create anything from CLIs to fully featured web apps.

The general philosophy is to provide loosely coupled libraries that can be composed as a suite of tools, vs. a `do it all` framework.

# Packages

The main packages are as follows:

- `assert` : helpers for writing tests; wraps `*testing.T` with more useful assertions.
- `collections` : common collections like ringbuffers and sets. 
- `configutil` : helpers for reading config files.
- `cron` : time triggered job management.
- `db` : our postgres orm.
- `db/migration` : helpers for writing postgres migrations.
- `env` : helpers for reading / writing / testing environment variables.
- `exception` : wraps error types with stack traces. 
- `logger` : our performance oriented event bus; event triggering is supported in most major packages.
- `oauth` : helpers for integrating with google oauth manager. 
- `proxy` : an http/https reverse proxy.
- `proxy/proxy` : a cli server the proxy.
- `request` : wrappers for `http.Client` with support for testing and a fluent api.
- `selector` : a portable implementation of kubernetes selectors.
- `semver` : semantic versioning helpers.
- `template` : text-template helpers.
- `template/template` : a cli for reading templates and outputting results.
- `util` : the junk drawer of random stuff. 
- `uuid` : generate and parse uuid v4's.
- `web` : our web framework; useful for both rest api's and view based apps.
- `workqueue` : a background work queue when you need to have a fixed number of workers.
- `yaml` : a yaml marshaller / unmarshaller. based on `go-yaml`.

# Code Style Notes

- Where possible, follow the [golang proverbs](https://go-proverbs.github.io/).
- The primary type a package exports should be creatable with a bare constructor `New()` unless there are non-trivial defaults to set.
- "Fluent APIs"
    - Mutators that return a reference to the receiver, and don't produce an error, should start with `With...()`.
    - This allows you to chain calls, ex. `New().WithFoo(...).WithBar(...)`.
    - Mutators that can return an error should start with `Set...()`
- Field accessors should be the uppercase name of the field, i.e. `foo` would have an accessor `Foo()`.
- Where possible, types should have a config object that fully represent the options you can set with `With` or `Set` mutators. 
- Said types should also have a constructor in the form `NewFromEnv` that uses the `go-sdk/env` package to read options set in the environment.
- Minimize dependencies between packages as much as possible; add external dependencies with *extreme* care.
    - our only current external dependencies are the golang stdlib and `github.com/lib/pq` for the `go-sdk/db`.

# Dependency Guidelines

- `assert` should depend only on the stdlib.
- `exception` should depend only on `assert` and the stdlib.
- `util` should depenend only on `exception`, `assert`, and the stdlib.
- `logger` should depend only on `util`, `exception`, `assert`, and the stdlib.
- Internal package dependencies otherwise are fair game, but try and minimize coupling.
- Do not add external packages unless absolutely necessary.
- If you do have to add an external dependency, make sure it's included in `make new-install`.

# Version Management

Generally we follow semantic versioning. What that means in practice:

> [major].[minor].[patch]

We increment the major version if there are *any* breaking changes. A breaking change is defined as a change that would cause code written against the current major version having a build failure of any type. That's even for a trivial find and replace. Once you merge a new api or name for an object after CR, that's it.

We increment minor versions if we add new things that don't cause breaking changes. 

We increment patch versions if we fix issues with the current set of objects.

## To increment the local version

Patch:
> make increment-patch

Minor:
> make increment-minor

Major:
> make increment-major