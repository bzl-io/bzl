load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library", "go_test")

go_library(
    name = "config",
    srcs = [
        "config.go",
    ],
    deps = [
        "@com_github_mitchellh_go_homedir//:go_default_library",                
    ],
    importpath = "github.com/bzl-io/bzl/config",
    visibility = ["//visibility:public"],
)

