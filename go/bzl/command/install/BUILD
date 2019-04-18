load("@io_bazel_rules_go//go:def.bzl", "go_library")

# go_library(
#     name = "install",
#     srcs = ["install.go"],
#     data = [
#     ],
#     importpath = "github.com/bzl-io/bzl/command/install",
#     visibility = ["//visibility:public"],
#     deps = [
#         "//go/bzl/command",
#         "//go/bzl/gh",
#         "@com_github_davecgh_go_spew//spew:go_default_library",
#         "@com_github_dustin_go_humanize//:go_default_library",
#         "@com_github_google_go_github//github:go_default_library",
#         "@com_github_joeybloggs_go_download//:go_default_library",
#         "@com_github_mitchellh_go_homedir//:go_default_library",
#         "@com_github_mitchellh_ioprogress//:go_default_library",
#         "@com_github_urfave_cli//:go_default_library",
#         "@org_golang_x_net//context:go_default_library",
#         #"@com_github_vbauerster_mpb//:go_default_library"
#         # These fail with cross-compilation
#         #"@com_github_cheggaaa_pb//:go_default_library",
#         #"@com_github_mattn_go_runewidth//:go_default_library",
#     ],
# )

go_library(
    name = "go_default_library",
    srcs = ["install.go"],
    importpath = "github.com/bzl-io/bzl/go/bzl/command/install",
    visibility = ["//visibility:public"],
    deps = [
        "@com_github_bzl_io_bzl//gh:go_default_library",
        "@com_github_dustin_go_humanize//:go_default_library",
        "@com_github_google_go_github//github:go_default_library",
        "@com_github_joeybloggs_go_download//:go_default_library",
        "@com_github_mitchellh_go_homedir//:go_default_library",
        "@com_github_mitchellh_ioprogress//:go_default_library",
        "@com_github_urfave_cli//:go_default_library",
        "@org_golang_x_net//context:go_default_library",
    ],
)
