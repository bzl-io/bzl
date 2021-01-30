---
layout: default
title: Rules
permalink: /rules
nav_order: 7
---

# Rules

The default tab shown when visiting a rule is the starlark `BUILD` rule formatted and labels linked to other nodes in the build graph:

<img width="720" alt="https://github.com/bzl-io/bzl/pull/12" src="https://user-images.githubusercontent.com/50580/106347195-ee174780-6279-11eb-8a91-50467b0a6338.png" style="border: 1px solid rgba(0,0,0,0.16)">

## Inputs and Outputs 

The **Io** tab displays the rule inputs and predeclared outputs:

<img width="720" alt="https://github.com/bzl-io/bzl/pull/12" src="https://user-images.githubusercontent.com/50580/106347233-3d5d7800-627a-11eb-9fa0-50ac5ec02b7c.png" style="border: 1px solid rgba(0,0,0,0.16)">

## Actions 

The **Actions** tab displays the rule actions that are spawned by the rule.  This is a view of `bazel aquery 'outputs(.*, LABEL)`:

<img width="720" alt="https://github.com/bzl-io/bzl/pull/12" src="https://user-images.githubusercontent.com/50580/106347798-5b2cdc00-627e-11eb-95e0-39566734f2ca.png" style="border: 1px solid rgba(0,0,0,0.16)">

## Attributes

The rule attributes documents the types and values of explicit and implicit attributes:

<img width="720" alt="https://github.com/bzl-io/bzl/pull/12" src="https://user-images.githubusercontent.com/50580/106348841-916e5980-6286-11eb-906b-7d5ac3a51bb3.png" style="border: 1px solid rgba(0,0,0,0.16)">

## Dependencies

The dependencies tab allows you to visualize the:

- reverse dependencies (what other targets depend on this rule within the default workspace)
- reverse dependencies within package (must faster to compute)
- rule dependencies (other rules this one depends on)

<img width="720" alt="https://github.com/bzl-io/bzl/pull/12" src="https://user-images.githubusercontent.com/50580/106348907-0e013800-6287-11eb-866e-c2904d2c2c0b.png" style="border: 1px solid rgba(0,0,0,0.16)">

## Events

The **Events** tab shows the build events following a build/test of this rule:

<img width="720" alt="https://github.com/bzl-io/bzl/pull/12" src="https://user-images.githubusercontent.com/50580/106349020-c6c77700-6287-11eb-81fc-2a5e4070d496.png" style="border: 1px solid rgba(0,0,0,0.16)">

This is covered in further detail in the [build events](/events) page.