# Rake Cloud Native Buildpack

## `gcr.io/paketo-community/rake`

The Rake CNB sets a command to run the default Rake task.

## Integration

This CNB writes a command, so there's currently no scenario we can
imagine that you would need to require it as dependency. If a user likes to
include some other functionality, it can be done independent of the Rake CNB
without requiring a dependency of it.

To package this buildpack for consumption:
```
$ ./scripts/package.sh
```
This builds the buildpack's source using GOOS=linux by default. You can supply another value as the first argument to package.sh.

## `buildpack.yml` Configurations

There are no extra configurations for this buildpack based on `buildpack.yml`.
