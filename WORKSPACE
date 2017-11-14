
# ================================================================

http_archive(
    name = "io_bazel_rules_go",
    url = "https://github.com/bazelbuild/rules_go/releases/download/0.7.0/rules_go-0.7.0.tar.gz",
    sha256 = "91fca9cf860a1476abdc185a5f675b641b60d3acf0596679a27b580af60bf19c",
)

load("@io_bazel_rules_go//go:def.bzl", "go_rules_dependencies", "go_register_toolchains", "go_repository")

load("@io_bazel_rules_go//proto:def.bzl", "proto_register_toolchains")

go_rules_dependencies()

go_register_toolchains()


proto_register_toolchains()

# ================================================================

go_repository(
    name = "com_github_google_go_github",
    importpath = "github.com/google/go-github",
    commit = "61bf8a6a2009b7fb88173ed454c77fb8013b8852",
)

go_repository(
    name = "com_github_google_go_querystring",
    importpath = "github.com/google/go-querystring",
    commit = "53e6ce116135b80d037921a7fdd5138cf32d7a8a",
)

go_repository(
    name = "com_github_davecgh_go_spew",
    importpath = "github.com/davecgh/go-spew",
    commit = "a476722483882dd40b8111f0eb64e1d7f43f56e4",
)

go_repository(
    name = "com_golang_google_genproto",
    importpath = "google.golang.org/genproto",
    commit = "3273178ea4684acc4f512f7bef7349dd72db88f6",
)

go_repository(
    name = "com_github_gregjones_httpcache",
    importpath = "github.com/gregjones/httpcache",
    commit = "22a0b1feae53974ed4cfe27bcce70dba061cc5fd",
)

go_repository(
    name = "com_github_peterbourgon_diskv",
    importpath = "github.com/peterbourgon/diskv",
    commit = "53ef9e43a0bc608e737e6bfed35207ad9cb1ad54",
)

go_repository(
    name = "com_github_google_btree",
    importpath = "github.com/google/btree",
    commit = "316fb6d3f031ae8f4d457c6c5186b9e3ded70435",
)

go_repository(
    name = "com_github_urfave_cli",
    importpath = "github.com/urfave/cli",
    commit = "44cb242eeb4d76cc813fdc69ba5c4b224677e799",
)

go_repository(
    name = "com_github_dustin_go_humanize",
    importpath = "github.com/dustin/go-humanize",
    commit = "6d15c0ae71e55ed645c21ac4945aaadbc0e9a590",
)

go_repository(
    name = "com_github_cheggaaa_pb",
    importpath = "github.com/cheggaaa/pb",
    commit = "657164d0228d6bebe316fdf725c69f131a50fb10",
)

go_repository(
    name = "com_github_mattn_go_runewidth",
    importpath = "github.com/mattn/go-runewidth",
    commit = "97311d9f7767e3d6f422ea06661bc2c7a19e8a5d",
)

go_repository(
    name = "com_google_bazelbuild_buildtools",
    importpath = "github.com/bazelbuild/buildtools",
    commit = "a861d1c9f86278f04ae7719ef48d22016781d766",
)

go_repository(
    name = "com_github_mitchellh_go_homedir",
    importpath = "github.com/mitchellh/go-homedir",
    commit = "b8bc1bf767474819792c23f32d8286a45736f1c6",
)

go_repository(
    name = "com_github_joeybloggs_go_download",
    importpath = "github.com/joeybloggs/go-download",
    commit = "26df310821f0e7614a736a2ee38ecd7e2f6ef6da",
)



#git_repository(
 # name = "org_pubref_rules_protobuf",
 # remote = "https://github.com/pubref/rules_protobuf",
 # tag = "v0.8.1",
    #)

#load("@org_pubref_rules_protobuf//go:rules.bzl", "go_proto_repositories")

#go_proto_repositories()
