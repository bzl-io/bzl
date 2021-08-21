---
layout: default
title: Remote Cache Overview
permalink: /remote-cache
nav_order: 3
---

# Remote Cache

Bzl has a built-in fast and lightweight gRPC remote cache implementation.  To
start the cache:

```sh
$ bzl cache
```

By default, `bzl` will create a filesystem directory at
`$USER_CACHE_DIR/bzl/remote-cache` and maintain the cache with a max-size of
10GB, evicting blobs on an as-needed basis using Least Recently Used (LRU)
semantics.

To use the cache, invoke bazel with the remote_cache flag as follows:
`--remote_cache=grpc://localhost:2020`.

Alternatively, put this on your ~/.bazelrc file:

```
build --remote_cache=grpc://localhost:2020
```

## Help

Optional CLI flags can be provided ti change the bind address, filesystem
directory, or maximum size of the cache:

```
bzl cache --help
start the remote disk cache

Usage:
  bzl cache [flags]

Flags:
      --address string    bind address for the remote cache. (default "grpc://localhost:2020")
      --dir string        base directory for the disk cache.
  -h, --help              help for cache
      --max_size_gb int   size in GB for the disk cache. (default 10)
```