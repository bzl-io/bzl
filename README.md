# bzl

`bzl` is a command-line application that makes it easier to work with multiple
versions of bazel.  You can install different versions and easily switch between
them.

## Install bzl

`bzl` ships as a single executable go binary. Download the file directly from
the [Github Releases Page](https://github.com/bzl-io/bzl/releases) for the
precompiled platform of your choice (or build from source).

## How is `bzl` pronounced?

`bzl` is pronounced like bezel, as in "*the bezel of a watch*". The name invokes
it's function (a wrapper around bazel).

> In the 1950s, watchmakers realized that an external bezel was the best way to
> add functions to a watch without complicating the movement, and so the
> external watch bezel was born.

## `bzl` commands

### `$ bzl --help`

Show help.

### `$ bzl install`

List or install available bazel installs.

Examples:

| Command | Description |
| --- | --- |
| `$ bzl install` | List all available releases |
| `$ bzl install 0.8.0` | Install bazel release 0.8.0 |
| `$ bzl install --list 0.8.0` | Show the assets bundled in install 0.8.0 |

### `$ bzl target`

Pretty-print available targets in the current workspace.

Example:

```
$ bzl targets
go_library        rule  //proto/bes:go_default_library
go_library        rule  //:go_default_library
go_library        rule  //command:go_default_library
go_library        rule  //command/targets:go_default_library
_gazelle_runner   rule  //:gazelle-runner
go_library        rule  //gh:go_default_library
go_proto_library  rule  //proto/bes:build_event_stream_go_proto
proto_library     rule  //proto/build:bzl_proto
proto_library     rule  //proto/bes:build_event_stream_proto
go_proto_library  rule  //proto/build:bzl_go_proto
go_library        rule  //command/release:go_default_library
sh_binary         rule  //:gazelle
go_library        rule  //command/install:go_default_library
go_library        rule  //config:go_default_library
go_test           rule  //:go_default_test
go_library        rule  //proto/build:go_default_library
go_library        rule  //bazelutil:go_default_library
go_binary         rule  //:bzl
go_library        rule  //command/use:go_default_library
_buildifier       rule  //:buildifier
```

### `$ bzl release`

Publish a release for a `go_binary` target.

This is a simple variant of the "goreleaser" tool for go binaries built by
rules_go.  Here's what the command does:

1. Reads the name of a `go_binary` target (example: `//:bzl`) as the first
   command line argument.
2. Takes a list of platform names via the `--platform` argument (example: `linux_amd64`)
   it runs the equivalent of `bazel build --platforms @io_bazel_rules_go//go/toolchain:linux_amd64`
3. Copies the output file to `{ASSET_DIR}/{LABEL_NAME}-{TAG}-{PLATFORM_NAME}`
   (example: `.assets/bzl-v0.1.3-linux-x64_64`) where `--asset_dir={ASSET_DIR}`,
   `--tag={TAG}`, and `--platform={PLATFORM}`.  The `{PLATFORM_NAME}` can be
   remapped to a different string via `--platform_name={PLATFORM}=NAME`
   (example: `windows_amd64=windows-x64_64`).
4. Reads a release notes file (example: `--notes=RELEASE.md`).
5. Uploads the staged assets to github.  Authentication is via environment
   variables (example: `BZL_GH_USERNAME=pcj`
   `BZL_GH_PASSWORD={PERSONAL_ACCESS_TOKEN}`)
6. Tags the release as `--tag={TAG}`. 

Example:

```
$ bzl release \
    --owner=bzl-io \
    --repo=bzl \
    --commit=91801a92ea21cd73471e5a83ad2519d1a3f257f0 \
    --tag=v0.1.3 \
    --notes=RELEASE.md \
    --platform=linux_amd64 \
    --platform=darwin_amd64 \
    --platform=windows_amd64 \
    //:bzl
```
