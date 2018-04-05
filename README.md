go-sdk
======

# Design Guidelines

- `assert` should reference nothing else outside the stdlib.
- `util` should reference assert and sub-packages only.
- `logger` should reference `util` and `assert` only.
- Everything else is fair game.
- Where possible, don't add external dependencies. Use the stdlib, or inline if small. 
    - If you have to add external dependencies, make sure they're in the `new-install` target in the makefile.
- Is it going to be multiple types / functions? If not, put it in util.

# Version Management

Generally we follow semantic versioning. What that means in practice:

> [major].[minor].[patch]

We increment the major version if there are *any* breaking changes. A breaking change is defined as a change that would cause code written against the current major version having a build failure of any type. That's even for a trivial find and replace. Once you merge a new api or name for an object after CR, that's it.

We increment minor versions if we add new things that don't cause breaking changes. 

We increment patch versions if we fix issues with the current set of objects.

## To increment the patch version

> make increment-patch