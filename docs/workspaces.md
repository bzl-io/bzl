---
layout: default
title: Workspaces / Packages
permalink: /workspaces
nav_order: 6
---

# Workspaces

## Default Workspace

By default the _default workspace_ is displayed as a table:

<img width="720" alt="https://github.com/bzl-io/bzl/pull/12" src="https://user-images.githubusercontent.com/50580/106345739-8fe56700-626f-11eb-928b-983c1be5ee2a.png" style="border: 1px solid rgba(0,0,0,0.16)">

You can also view it as a tree (button at upper right):

<img width="720" alt="https://github.com/bzl-io/bzl/pull/12" src="https://user-images.githubusercontent.com/50580/106345823-2a45aa80-6270-11eb-8eaf-eb44aafad682.png" style="border: 1px solid rgba(0,0,0,0.16)">

Other buttons:

- **[All Rules]**: show all rules defined in the workspace.  *Caution: this can trigger a large bazel query and take a while to complete*.
- **[Repository Rule]**: show the repository rule where the workspace is declared.  Relevant for external workspaces.
- **[Open]**: open the workspace in your IDE.
- **[Run]**: build/test/run arbitrary targets command in the workspace.
- **[Build All]**: build all targets in the workspace.
- **[Test All]**: test all targets in the workspace.

## External Workspaces

Select the **External** tab to view the list of external workspaces.  This is similar to running `bazel query //external:*`:

<img width="720" alt="https://github.com/bzl-io/bzl/pull/12" src="https://user-images.githubusercontent.com/50580/106346032-4bf36180-6271-11eb-921e-237e177651ff.png" style="border: 1px solid rgba(0,0,0,0.16)">

The table view restricts the list to ones not declared in `@bazel_tools//...`.  If you'd like to view the complete list, use the tree view:

<img width="720" alt="https://github.com/bzl-io/bzl/pull/12" src="https://user-images.githubusercontent.com/50580/106346057-8b21b280-6271-11eb-9e54-3e91164bcc39.png" style="border: 1px solid rgba(0,0,0,0.16)">

