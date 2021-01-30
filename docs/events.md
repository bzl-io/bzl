---
layout: default
title: Events
permalink: /events
nav_order: 8
---

# Events

Bazel has a fairly sophisticated protocol for reporting events that occur during
the progress of a build.  Bzl acts as a `bes_backend` and typically runs via the
following flags:

```sh 
--bes_backend=localhost:1080
--bes_timeout=5s
--build_event_publish_all_actions=true
```

Events are streamed back up into the browser and presented in a timeline view
that allows you to see:

- all events by time
- all events of a particular type
- specific time window

To view events in realtime during a build or test invocation, click the
**[Build]** or **[Test]** buttons:

<img width="760" alt="https://github.com/bzl-io/bzl/pull/12" src="https://user-images.githubusercontent.com/50580/106349273-042d0400-628a-11eb-8b0c-bef6b07b65d6.gif">

Dragging a window/level will filter the event within that time window:

<img width="760" alt="https://github.com/bzl-io/bzl/pull/12" src="https://user-images.githubusercontent.com/50580/106349368-b238ae00-628a-11eb-8210-0b40d13ec090.gif">

You can clear the selection by clicking away, or use the menu:

<img width="760" alt="https://github.com/bzl-io/bzl/pull/12" src="https://user-images.githubusercontent.com/50580/106349441-0e9bcd80-628b-11eb-83ad-5759400ca9a9.gif">

Click on an event type to filter to only that type:


<img width="760" alt="https://github.com/bzl-io/bzl/pull/12" src="https://user-images.githubusercontent.com/50580/106349497-70f4ce00-628b-11eb-8a85-a16eb9516eb6.png">


## Event Types

A brief explanation of the types of build events (not an exhaustive list):

1. **Started**: provides a summary of the build metadata including start time,
   the version of bazel, the command, description of options.  This is a good
   overview.
2. **Structured Command Line**: an exhaustive list of the precise configuration
   provided by the bazel client.  Use this to see the exact set of flags used
   during the build.
3. **Options Parsed**: another view of the options.
4. **Build Metadata**: user configurable metadata about the build.  Can be
   populated with the `--build_metadata` flag (e.g.
   `--build_metdata=JOB_ID=12345`).
5. **Progress**: console logging, typically delivered via **stderr**.
6. **Workspace Configuration**: provides the execution root.
7. **Configuration**: the configuration for the build including such things as
   `cpu` and `platform`.
8. **Configured**: the rule and label after the configuration has been resolved.
9. **Workspace Status**: key value pairs populated via the `--workspace_status_command`.
10. **Named Set of Files**: a set of files.  In order to be efficient, bazel groups files into sets to avoid duplication where possible.
11. **Action**: the individual units of work involved in a build.  Each one
    represents a process that was executed either on the local workstation or a
    remote machine.  Each action has a `mnemonic`, an *action key* hash, a
    command line to execute, and a list of environment variables.  The mnemonic
    is sort of like a type name that works with flags such as `--spawn_strategy`
    and others.
12. **Test Result**: an event that reports test success or failure, whether it
    was cached, and links to the log file.
13. **Build Metrics**: a summary of actions executed, memory usage, etc.
14. **Build Tool Logs**: a summary of the time spent including the critical path.
15. **Completed**: reports files outputted by a target in their respective
    output groups.
16. **Aborted**: if the build was aborted.
17. **Finished**: overall success or failure.

These represent the inner *Bazel Build Events*.  Each one is wrapped in an outer
*Ordered Build Event* envelope that represents a more generic event container
that imposes a sequencing mechanism.
