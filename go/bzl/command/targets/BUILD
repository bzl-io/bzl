load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "targets",
    srcs = [
        "targets.go",
    ],
    deps = [
        "//go/bzl/bazel",
        "//go/bzl/command",
        "//proto:build_go",
        "@com_github_urfave_cli//:go_default_library",
        "@com_github_davecgh_go_spew//spew:go_default_library",
    ],
    importpath = "github.com/bzl-io/bzl/command/targets",
    visibility = ["//visibility:public"],
    data = [
    ],
)
