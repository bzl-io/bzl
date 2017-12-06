load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library", "go_test")

go_library(
    name = "bazel",
    srcs = [
        "bazel.go",
    ],
    deps = [
        "//proto:build_go",
        "//proto:build_event_stream_go",
        "//go/bzl/config",
        "@com_github_golang_protobuf//ptypes/any:go_default_library",
        "@com_github_matttproud_golang_protobuf_extensions//pbutil:go_default_library",        
        "@com_github_golang_protobuf//proto:go_default_library",
    ],
    importpath = "github.com/bzl-io/bzl/bazel",
    visibility = ["//visibility:public"],
)

