load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "release",
    srcs = [
        "release.go",
    ],
    deps = [
        "//go/bzl/bazel",
        "//go/bzl/gh",
        "//go/bzl/command",
        "//proto:build_event_stream_go",
        "@org_golang_x_net//context:go_default_library",
        "@org_golang_x_sync//errgroup:go_default_library",
        "@com_github_google_go_github//github:go_default_library",        
        "@com_github_urfave_cli//:go_default_library",
        "@com_github_pkg_errors//:go_default_library",
        "@com_github_davecgh_go_spew//spew:go_default_library",        
    ],
    importpath = "github.com/bzl-io/bzl/command/release",
    visibility = ["//visibility:public"],
    data = [
    ],
)
