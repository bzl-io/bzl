load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "install",
    srcs = [
        "install.go",
    ],
    deps = [
        "//go/bzl/command",
        "//go/bzl/gh",
        "@org_golang_x_net//context:go_default_library",
        "@com_github_google_go_github//github:go_default_library",        
        "@com_github_urfave_cli//:go_default_library",
        "@com_github_davecgh_go_spew//spew:go_default_library",        
        "@com_github_joeybloggs_go_download//:go_default_library",
        "@com_github_mitchellh_go_homedir//:go_default_library",        
        "@com_github_mitchellh_ioprogress//:go_default_library",        
        "@com_github_dustin_go_humanize//:go_default_library",
        #"@com_github_vbauerster_mpb//:go_default_library"
        # These fail with cross-compilation
        #"@com_github_cheggaaa_pb//:go_default_library",
        #"@com_github_mattn_go_runewidth//:go_default_library",
    ],
    importpath = "github.com/bzl-io/bzl/command/install",
    visibility = ["//visibility:public"],
    data = [
    ],
)
