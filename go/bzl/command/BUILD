load("@io_bazel_rules_go//go:def.bzl", "go_library")

# go_library(
#     name = "command",
#     srcs = ["command.go"],
#     data = [
#     ],
#     importpath = "github.com/bzl-io/bzl/command",
#     visibility = ["//visibility:public"],
#     deps = [
#         "@com_github_urfave_cli//:go_default_library",
#     ],
# )

go_library(
    name = "go_default_library",
    srcs = ["command.go"],
    importpath = "github.com/bzl-io/bzl/go/bzl/command",
    visibility = ["//visibility:public"],
    deps = ["@com_github_urfave_cli//:go_default_library"],
)
