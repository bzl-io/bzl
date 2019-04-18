load("@io_bazel_rules_go//go:def.bzl", "go_library")

# go_library(
#     name = "gh",
#     srcs = ["gh.go"],
#     importpath = "github.com/bzl-io/bzl/gh",
#     visibility = ["//visibility:public"],
#     deps = [
#         "@com_github_google_go_github//github:go_default_library",
#         "@com_github_google_go_querystring//query:go_default_library",
#         "@com_github_gregjones_httpcache//:go_default_library",
#         "@com_github_gregjones_httpcache//diskcache:go_default_library",
#         "@com_github_mitchellh_go_homedir//:go_default_library",
#         #"@org_golang_x_oauth2//internal:go_default_library",
#         #"@org_golang_x_oauth2//:go_default_library",
#     ],
# )

go_library(
    name = "go_default_library",
    srcs = ["gh.go"],
    importpath = "github.com/bzl-io/bzl/go/bzl/gh",
    visibility = ["//visibility:public"],
    deps = [
        "@com_github_google_go_github//github:go_default_library",
        "@com_github_gregjones_httpcache//:go_default_library",
        "@com_github_gregjones_httpcache//diskcache:go_default_library",
        "@com_github_mitchellh_go_homedir//:go_default_library",
    ],
)
