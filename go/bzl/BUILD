load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

# go_library(
#     name = "app",
#     srcs = ["app.go"],
#     importpath = "github.com/bzl-io/bzl",
#     visibility = ["//visibility:public"],
#     deps = [
#         "//go/bzl/bazel",
#         "//go/bzl/command/install",
#         "//go/bzl/command/targets",
#         "@com_github_urfave_cli//:go_default_library",
#     ],
# )

# sh_test(
#     name = "install_test",
#     srcs = [
#         "install_test.sh",
#     ],
#     deps = [
#     ],
#     data = [
#         ":bzl",
#     ],
# )

go_library(
    name = "go_default_library",
    srcs = ["app.go"],
    importpath = "github.com/bzl-io/bzl/go/bzl",
    visibility = ["//visibility:public"],
    deps = [
        "@com_github_bzl_io_go//bzl/bazel:go_default_library",
        "@com_github_bzl_io_go//bzl/command/install:go_default_library",
        "@com_github_bzl_io_go//bzl/command/targets:go_default_library",
        "@com_github_urfave_cli//:go_default_library",
    ],
)

go_test(
    name = "go_default_test",
    srcs = ["app_test.go"],
    embed = [":go_default_library"],
    deps = ["//:go_default_library"],
)
