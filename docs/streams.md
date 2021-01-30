---
layout: default
title: Streams
permalink: /streams
nav_order: 10
---

# Build Event Streams

Bzl acts as a *Build Event Stream Backend*, so you can point your bazel builds to Bzl and view the build events.

## Invoking a Stream

Example:

```sh 
bazel build //absl/base 
    --bes_backend=grpc://localhost:1080 
    --bes_results_url=http://localhost:8080/stream
```


~~~sh
INFO: Streaming build results to: http://localhost:8080/stream/d948b5fb-b9da-493b-8b14-a9af2ad076aa
INFO: Analyzed target //absl/base:base (0 packages loaded, 0 targets configured).
INFO: Found 1 target...
[0 / 1] [Prepa] BazelWorkspaceStatusAction stable-status.txt
Target //absl/base:base up-to-date:
  bazel-bin/absl/base/libbase.a
  bazel-bin/absl/base/libbase.so
INFO: Elapsed time: 0.182s, Critical Path: 0.00s
INFO: Build completed successfully, 1 total action
INFO: Streaming build results to: http://localhost:8080/stream/d948b5fb-b9da-493b-8b14-a9af2ad076aa
~~~

## Viewing

Click on the link at the end of the build to view the invocation details:

<img width="720" alt="https://github.com/bzl-io/bzl/pull/12" src="https://user-images.githubusercontent.com/50580/106352660-2599ea00-62a2-11eb-8e80-799f6b85d381.png" style="border: 1px solid rgba(0,0,0,0.16)">
