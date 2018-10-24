go-sdk
======

![Build Status](https://circleci.com/gh/blend/go-sdk.svg?style=shield&circle-token=:circle-token)

`go-sdk` is our core library of packages. These packages can be composed to create anything from CLIs to fully featured web apps.

The general philosophy is to provide loosely coupled libraries that can be composed as a suite of tools, vs. a `do it all` framework.

# Addtional CLI Tools

We also provide the following CLI tools to help with development that leverage some of these packages:

- `cmd/cover` : allows for project level coverage reporting and enforcement.
- `cmd/profanity` : profanity rules checking (i.e. fail on grep match).
- `cmd/recover` : recover crashed processes (to be used when debugging panics).
- `cmd/semver` : semver maniuplation and validation. 

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
- Anything that can return an error, should. Anything that needs to return a single value (but would return an error) should panic on that error and should be prefixed by `Must...`.
- Minimize dependencies between packages as much as possible; add external dependencies with *extreme* care.

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
