go-sdk
======

[![Build Status](https://circleci.com/gh/blend/go-sdk.svg?style=shield&circle-token=:circle-token)

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

# Design Guidelines

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

## To increment the patch version

> make increment-patch