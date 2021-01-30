---
layout: default
title: UI Overview
permalink: /ui
nav_order: 3
---

# User Interface

## Content Areas

<img width="720" alt="https://github.com/bzl-io/bzl/pull/12" src="https://user-images.githubusercontent.com/50580/106352334-765c1380-629f-11eb-9e88-7df13518bf67.png" style="border: 1px solid rgba(0,0,0,0.16)">

The content area is divided into three main regions:

1. [Repositories](/bzl/repositories): the bazel workspaces on your workstation.
2. [Streams](/bzl/streams): a view of build event protocol invocations.
3. [Community](/bzl/community): a list of well-known bazel repositories waiting to be explored.

## Menus

The menu system is context sensitive, meaning the menus change according to the component tree being visited.

<img width="720" alt="https://github.com/bzl-io/bzl/pull/12" src="https://user-images.githubusercontent.com/50580/106351772-79550500-629b-11eb-8b16-a186328e710d.png" style="border: 1px solid rgba(0,0,0,0.16)">

## Keyboard Shortcuts

Keyboard shortcuts are listed in the menus. They are generally 2-key sequences such as `r`, `a` (rule attributes).

Single letter keys such as `b` (build), `t` (test) are reserved for common operations.  Application home is mapped to the backtick '`'

The search bar is focused via the forward slash key `/`.

## Search

Search is populated incrementally via the things (mostly targets) that have been
visited thus far in your browsing session.

<img width="720" alt="https://github.com/bzl-io/bzl/pull/12" src="https://user-images.githubusercontent.com/50580/106351868-292a7280-629c-11eb-838e-6d0923f5f056.gif">

Unfortunately this does not represent a comprehensive list of targets in the
*universe* of the build.  The reasoning is that populating a target index is an
expensive operation and is undesirable to lock the bazel server for this
purpose.

A future improvement may address incrementally updating a target database.

## Shortcomings 

The UI is a single page application but still may require page reload at times.
Bzl may trigger bazel queries that take a long time to complete and can be
difficult to cancel.
