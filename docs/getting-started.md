---
layout: default
title: Getting Started
permalink: /getting-started
nav_order: 1
---

# Getting Started

## Registration

Go to <https://build.bzl.io> and [login](https://build.bzl.io/bzl/login) via GitHub or Email/Password.

Once logged in, choose a subscription plan and enter your card details:

<img width="631" alt="https://github.com/bzl-io/bzl/pull/12" src="https://user-images.githubusercontent.com/50580/106341648-45f28600-625b-11eb-8c5b-fcc1e513d620.png" style="border: 1px solid rgba(0,0,0,0.16)">

Once confirmed, you'll be forwarded to the download page:

<img width="625" alt="Screen Shot 2021-01-29 at 6 05 45 PM" src="https://user-images.githubusercontent.com/50580/106342048-a635f780-625c-11eb-9779-d60a18436c4a.png" style="border: 1px solid rgba(0,0,0,0.16)">

Copy the `curl` command and paste into your shell.  Copy the auth token when prompted.

The binary will be installed at `~/.bzl/bin/bzl`; the `~/.bzl/bin` will be added to your `PATH`.

<img width="625" alt="Screen Shot 2021-01-29 at 6 14 29 PM" src="https://user-images.githubusercontent.com/50580/106342390-e6e24080-625d-11eb-82f3-1265eb062b2b.png" style="border: 1px solid rgba(0,0,0,0.16)">

## Usage

Here's brief summary of the CLI operations.

Any bazel-native command is passed through directly to the underlying bazel:

```sh
$ bzl build //foo:bar
$ bzl query deps(//foo:bar) --output label_kind
...
```

### serve/open

Open the user interface in the current workspace:

```sh
$ bzl serve           # starts webserver
$ bzl open            # opens browser tab in current workspace
$ bzl open //foo:bar  # opens browser tab at given label
```

### license 

Print or renew your license:

```
$ bzl license info
$ bzl license renew
Backed-up previous license to /Users/foo/.bzl/license.key.Updated /Users/foo/.bzl/license.key
```

### use

Use is a handy repository rule generator.  Never write another repository rule
by hand again! (or, at least less frequently).  Examples:

```sh 
$ bzl use rules_proto
````

```python
# Branch: master
# Commit: a0761ed101b939e19d83b2da5f59034bffc19c12
# Date: 2021-01-26 15:30:54 +0000 UTC
# URL: https://github.com/bazelbuild/rules_proto/commit/a0761ed101b939e19d83b2da5f59034bffc19c12
#
# Merge pull request #81 from Yannic/patch-3
#
# Bump bazel-toolchains to 3.7.2
# Size: 11622 (12 kB)
http_archive(
    name = "rules_proto",
    sha256 = "2a20fd8af3cad3fbab9fd3aec4a137621e0c31f858af213a7ae0f997723fc4a9",
    strip_prefix = "rules_proto-a0761ed101b939e19d83b2da5f59034bffc19c12",
    urls = ["https://github.com/bazelbuild/rules_proto/archive/a0761ed101b939e19d83b2da5f59034bffc19c12.tar.gz"],
)
```

```sh 
$ bzl use go github.com/google/uuid v1.2.0
```

```sh 
# Release: v1.2.0
# TargetCommitish: master
# Date: 2021-01-22 18:20:15 +0000 UTC
# URL: https://github.com/google/uuid/releases/tag/v1.2.0
# Size: 14158 (14 kB)
go_repository(
    name = "com_github_google_uuid",
    importpath = "github.com/google/uuid",
    tag = "v1.2.0",
)
```

```
$ bzl use file https://cdnjs.cloudflare.com/ajax/libs/underscore.js/1.12.0/underscore-min.js
```

```py 
# HTTP/2.0 200 OK
# Expires: Thu, 20 Jan 2022 02:01:30 GMT
# Last-Modified: Mon, 21 Dec 2020 09:19:03 GMT
# Server: cloudflare
# Size: 19358 (19 kB)
http_file(
    name = "cdnjs_cloudflare_com_ajax_libs_underscore_js_1_12_0_underscore_min_js",
    sha256 = "1bc0ea4e2fe66ac337fb1863bbdb4c8f044ee4e84dbe0f0f1b3959bebfa539c1",
    urls = ["https://cdnjs.cloudflare.com/ajax/libs/underscore.js/1.12.0/underscore-min.js"],
)
```

### install

List or view published bazel releases:

```sh 
$ bzl install
4.0.0
3.7.2 (installed)
3.7.1 (installed)
3.7.0 (installed)
3.6.0 (installed)
...

$ bzl install 4.0.0
2021/01/29 18:56:56 Downloading https://releases.bazel.build/4.0.0/release/bazel-4.0.0-darwin-x86_64...
```

> This caches downloaded binaries in the same location as `bazelisk`.

Bzl also functions as an HTTP or gRPC remote cache:

```sh
$ bazel cache
Bazel Cache: /Users/i868039/Library/Caches/bzl/bazel-remote-cache
  --remote_cache=http://localhost:2080
  --remote_cache=grpc://localhost:2092
```

The language server is typically only used by the VSCode extension, but it
functions as an LSP server for `BUILD` files.

```sh 
$ bzl lsp serve
```

## Support

Please contact <hello@stack.build> for support for file an issue at <https://github.com/bzl-io/bzl/issues>.