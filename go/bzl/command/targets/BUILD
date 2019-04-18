load("@io_bazel_rules_go//go:def.bzl", "go_library")

# go_library(
#     name = "targets",
#     srcs = ["targets.go"],
#     data = [
#     ],
#     importpath = "github.com/bzl-io/bzl/command/targets",
#     visibility = ["//visibility:public"],
#     deps = [
#         "//go/bzl/bazel",
#         "//go/bzl/command",
#         "//proto:build_go",
#         "@com_github_davecgh_go_spew//spew:go_default_library",
#         "@com_github_urfave_cli//:go_default_library",
#     ],
# )

go_library(
    name = "go_default_library",
    srcs = ["targets.go"],
    importpath = "github.com/bzl-io/bzl/go/bzl/command/targets",
    visibility = ["//visibility:public"],
    deps = [
        "@com_github_bzl_io_bzl//bazel:go_default_library",
        "@com_github_bzl_io_bzl//proto/build_go:go_default_library",
        "@com_github_urfave_cli//:go_default_library",
    ],
)
