---
layout: default
title: Repositories
permalink: /repositories
nav_order: 5
---

# Repositories

Bzl uses the terminology _repository_ to refer to a collection of files rooted at the `WORKSPACE` file.  This includes the _default workspace_ (e.g., anything under `@//...`) as well as the entire set of _external workspaces_ (e.g. `@foo/...`).

To view the list of repositories on your workstation visit <http://127.0.0.1:8080/local>:

<img width="720" alt="https://github.com/bzl-io/bzl/pull/12" src="https://user-images.githubusercontent.com/50580/106345046-84dc0800-626a-11eb-9019-557a3df8f5a4.png" style="border: 1px solid rgba(0,0,0,0.16)">

### Repository Discovery

Bzl discovers bazel repositories on your workstation by inspecting the OS
process list.  From there it scans the `--output_user_root`(s) named by the
bazel server process.

If you'd like to hardcode the location of repos on your filesystem, you can put
the absolute path to the repository in your `~/.bzlrc` file:

```sh 
common --repository_dir=/path/to/foo
common --repository_dir=/path/to/bar
common --repository_dir=/path/to/baz
```

### Abandoned Output Bases

If bzl discovers a bazel _output base_ (the place where it stores cached data)
that has no corresponding workspace, it labels this as an _Abandoned Output
Base_.

An output base can be abandoned when the filesystem directory where the
`WORKSPACE` lived is deleted, but a `bazel clean --expunge` was not issued
first.  It can also occur when a WORKSPACE is renamed.

<img width="720" alt="https://github.com/bzl-io/bzl/pull/12" src="https://user-images.githubusercontent.com/50580/106345508-f79ab280-626d-11eb-94c4-8263dd9dda25.png" style="border: 1px solid rgba(0,0,0,0.16)">

Over time, this can consume significant disk space.

NOTE: in case you were wondering, you can nuke the entire `--output_user_root`
to reclaim all disk space being used by bazel. For example, the following is
completely safe from a data-loss perspective:

```sh
$ rm -rf /var/tmp/_bazel_foo/
```