load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library", "go_test")

go_library(
    name = "app",
    srcs = [
        "app.go",
    ],
    deps = [
        "//go/bzl/command/install",
        "//go/bzl/command/targets",
        "@com_github_urfave_cli//:go_default_library",
    ],
    importpath = "github.com/bzl-io/bzl",
    visibility = ["//visibility:public"],
)

go_binary(
    name = "bzl",
    srcs = [
        "bzl.go",
    ],
    deps = [
        ":app",
    ],
    visibility = ["//visibility:public"],
    data = [
    ],
)

