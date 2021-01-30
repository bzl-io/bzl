---
layout: default
title: Configuration
permalink: /configuration
nav_order: 4
---

# Configuration

Bzl is configured in a manner similar to bazel itself: either by command line
flags or an rcfile(s).

## Command Line flags

Refer to the CLI documentation for a complete reference.  For example:

```sh
$ bzl serve --address=":4567"   # change the host/port(s) for network binding
```

## Configration file ~/.bzlrc

The optional file `${HOME}/.bzlrc` can be used to configure flags.  The file
format is similar to bazel itself having the form `COMMAND --list --of --flags`.

```sh
# Change the base directory where bzl caches files
common --base_dir=/tmp/bzl
```

## Editor Integration

The UI allows you to open files in your preferred editor/IDE.  By default, this
is configured for vscode.  To adjust this, please change the following in your
`~/.bzlrc` file:

```sh
common --open_command='code --goto {FILE}:{LINE}:{COLUMN}'  # vscode
common --open_command='gvim +{LINE} {FILE}'                 # vi
common --open_command='emacsclient +{LINE} {FILE}'          # emacs
common --open_command='idea {FILE} --line {LINE}'           # intellij
common --open_command='atom {FILE}:{LINE}:{COLUMN}'         # atom
common --open_command='subl {FILE}:{LINE}:{COLUMN}'         # sublime
```

- The string `{FILE}` will be replaced with the actual filename or dirname.  
- The string `{LINE}` will be replaced with the desired line number.
- The string `{COLUMN}` will be replaced with the desired column number.

> Note that *this only works with an editor command that does not block* (for
> example, 'vi').


## GitHub Integration

As various times `bzl` makes calls to the GitHub API.  You can experience API
request limits without additional configuration of your github credentials.  A
recommended configuration to address this:

```bash
export GITHUB_USERNAME=<PERSONAL_ACCESS_TOKEN>
export GITHUB_PASSWORD=x-oauth-basic
```

> Visit <https://github.com/settings/tokens> to allocate a personal access token.

## Logging / Debugging

You can set the environment variable `LOG_LEVEL` to `warn`, `info`, or `debug`
to get additional contextual logging.

More selective "per-service" logging can be enabled via the syntax
`LOG_LEVEL={SERVICE_NAME}=debug`.  The service names canbe discovered by running
with `LOG_LEVEL=debug`.
