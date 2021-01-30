---
layout: default
title: Codesearch
permalink: /codesearch
nav_order: 11
---

# Codesearch

Bzl implements a codesearch feature that allows you to define *codesearch indexes* based on bazel queries, and search them using simple or regular expressions.

## Creating an Index

To create an index, provide the bazel query as the positional arguments and give it a name:

```sh 
$ bzl code index create --name=all 'deps(//absl/...)'
```

You can search the index via the command line:

```sh
bzl code search --index=all '<sstream>'
absl/random/internal/nonsecure_base_test.cc:21: #include <sstream>
embedded_tools/src/main/cpp/util/errors_windows.cc:20: #include <sstream>
absl/random/examples_test.cc:17: #include <sstream>
/private/var/tmp/_bazel_i868039/9a22f63cfea7c4a7c8ae084f584bea24/external/com_github_google_benchmark/src/sysinfo.cc:59: #include <sstream>
embedded_tools/src/main/cpp/util/errors_posix.cc:17: #include <sstream>
embedded_tools/src/main/native/windows/file.cc:26: #include <sstream>
embedded_tools/src/tools/launcher/util/launcher_util.cc:30: #include <sstream>
absl/strings/internal/str_format/bind.cc:19: #include <sstream>
```

Or via the UI:

Use the **Repository > Search Code** menu option (`o` `k`) to navigate to the codesearch widget:

<img width="720" alt="https://github.com/bzl-io/bzl/pull/12" src="https://user-images.githubusercontent.com/50580/106353568-c095c280-62a8-11eb-9b98-e0a4f4484db9.gif">
