load("@io_bazel_rules_go//go:def.bzl", "go_library")

# go_library(
#     name = "bazel",
#     srcs = ["bazel.go"],
#     importpath = "github.com/bzl-io/bzl/bazel",
#     visibility = ["//visibility:public"],
#     deps = [
#         "//go/bzl/config",
#         "//proto:build_event_stream_go",
#         "//proto:build_go",
#         "@com_github_golang_protobuf//proto:go_default_library",
#         "@com_github_golang_protobuf//ptypes/any:go_default_library",
#         "@com_github_matttproud_golang_protobuf_extensions//pbutil:go_default_library",
#     ],
# )

go_library(
    name = "go_default_library",
    srcs = ["bazel.go"],
    importpath = "github.com/bzl-io/bzl/go/bzl/bazel",
    visibility = ["//visibility:public"],
    deps = [
        "@com_github_bzl_io_bzl//proto/build_event_stream_go:go_default_library",
        "@com_github_bzl_io_bzl//proto/build_go:go_default_library",
        "@com_github_bzl_io_go//bzl/config:go_default_library",
        "@com_github_golang_protobuf//proto:go_default_library",
        "@com_github_matttproud_golang_protobuf_extensions//pbutil:go_default_library",
    ],
)
