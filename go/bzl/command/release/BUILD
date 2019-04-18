load("@io_bazel_rules_go//go:def.bzl", "go_library")

# go_library(
#     name = "release",
#     srcs = ["release.go"],
#     data = [
#     ],
#     importpath = "github.com/bzl-io/bzl/command/release",
#     visibility = ["//visibility:public"],
#     deps = [
#         "//go/bzl/bazel",
#         "//go/bzl/command",
#         "//go/bzl/gh",
#         "//proto:build_event_stream_go",
#         "@com_github_davecgh_go_spew//spew:go_default_library",
#         "@com_github_google_go_github//github:go_default_library",
#         "@com_github_pkg_errors//:go_default_library",
#         "@com_github_urfave_cli//:go_default_library",
#         "@org_golang_x_net//context:go_default_library",
#         "@org_golang_x_sync//errgroup:go_default_library",
#     ],
# )

go_library(
    name = "go_default_library",
    srcs = ["release.go"],
    importpath = "github.com/bzl-io/bzl/go/bzl/command/release",
    visibility = ["//visibility:public"],
    deps = [
        "@com_github_bzl_io_bzl//bazel:go_default_library",
        "@com_github_bzl_io_bzl//gh:go_default_library",
        "@com_github_bzl_io_bzl//proto/build_event_stream_go:go_default_library",
        "@com_github_davecgh_go_spew//spew:go_default_library",
        "@com_github_google_go_github//github:go_default_library",
        "@com_github_pkg_errors//:go_default_library",
        "@com_github_urfave_cli//:go_default_library",
        "@org_golang_x_net//context:go_default_library",
        "@org_golang_x_sync//errgroup:go_default_library",
    ],
)
