load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library", "go_test")

go_library(
    name = "app",
    srcs = [
        "app.go",
    ],
    deps = [
        "//go/bzl/bazel",
        "//go/bzl/plugin",
        "//go/bzl/command/install",
        "//go/bzl/command/targets",
        "//go/bzl/command/release",
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
    linkstamp = "github.com/bzl-io/bzl",
    x_defs = {
        "github.com/bzl-io/bzl.BUILD_SCM_REVISION": "{BUILD_SCM_REVISION}",
    },
)

sh_test(
    name = "install_test",
    srcs = [
        "install_test.sh",
    ],
    deps = [
    ],
    data = [
        ":bzl",
    ],
)

