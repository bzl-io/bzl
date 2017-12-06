load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "plugin",
    srcs = [
        "plugin.go",
    ],
    deps = [
        "//proto:api_go",
        "//go/bzl/config",
        "//go/bzl/command",
        "@org_golang_x_net//context:go_default_library",
        "@org_golang_x_net//trace:go_default_library",
        "@com_github_golang_glog//:go_default_library",
        "@org_golang_google_grpc//:go_default_library",
        "@org_golang_google_grpc//codes:go_default_library",
        "@org_golang_google_grpc//grpclog:go_default_library",
        "@org_golang_google_grpc//metadata:go_default_library",
    ],
    importpath = "github.com/bzl-io/bzl/plugin",
    visibility = ["//visibility:public"],
)
