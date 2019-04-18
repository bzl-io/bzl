load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive", "http_file")

local_repository(
    name = "build_stack_rules_proto",
    path = "/home/pcj/go/src/github.com/stackb/rules_proto",
)

load("@build_stack_rules_proto//:deps.bzl", "io_bazel_rules_go", "bazel_gazelle", "com_github_bazelbuild_buildtools")

io_bazel_rules_go()

bazel_gazelle()

com_github_bazelbuild_buildtools()

load("@io_bazel_rules_go//go:def.bzl", "go_rules_dependencies", "go_register_toolchains")

go_rules_dependencies()

go_register_toolchains()

load("@bazel_gazelle//:deps.bzl", "gazelle_dependencies", "go_repository")

gazelle_dependencies()

# gazelle:repo bazel_gazelle

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
    name = "com_github_pkg_errors",
    importpath = "github.com/pkg/errors",
    commit = "f15c970de5b76fac0b59abb32d62c17cc7bed265",
)

# go_repository(
#     name = "com_github_mattn_go_runewidth",
#     importpath = "github.com/mattn/go-runewidth",
#     commit = "97311d9f7767e3d6f422ea06661bc2c7a19e8a5d",
# )

go_repository(
    name = "org_golang_x_tools",
    importpath = "github.com/golang/tools",
    commit = "9b61fcc4c548d69663d915801fc4b42a43b6cd9c",
)

go_repository(
    name = "org_golang_x_sync",
    importpath = "github.com/golang/sync",
    commit = "fd80eb99c8f653c847d294a001bdf2a3a6f768f5",
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

# There exist multiple options for terminal-based progress bars, but
# this is the only one I could find that cross-compiles cleanly with
# current version of rules_go.  Not super fancy though.
go_repository(
    name = "com_github_mitchellh_ioprogress",
    importpath = "github.com/mitchellh/ioprogress",
    commit = "8163955264568045f462ae7e2d6d07b2001fc997",
)

go_repository(
    name = "com_github_joeybloggs_go_download",
    importpath = "github.com/joeybloggs/go-download",
    commit = "26df310821f0e7614a736a2ee38ecd7e2f6ef6da",
)

go_repository(
    name = "com_github_rs_cors",
    importpath = "github.com/rs/cors",
    commit = "eabcc6af4bbe5ad3a949d36450326a2b0b9894b8",  # Aug 1
)

go_repository(
    name = "com_github_gorilla_mux",
    importpath = "github.com/gorilla/mux",
    commit = "3f19343c7d9ce75569b952758bd236af94956061",
)

go_repository(
    name = "com_github_gorilla_context",
    importpath = "github.com/gorilla/context",
    commit = "08b5f424b9271eedf6f9f0ce86cb9396ed337a42",
)

go_repository(
    name = "com_github_vbauerster_mpb",
    importpath = "github.com/vbauerster/mpb",
    commit = "d3da256ab98a80013319df621e48db638748c044",
)

go_repository(
    name = "com_github_matttproud_golang_protobuf_extensions",
    importpath = "github.com/matttproud/golang_protobuf_extensions",
    commit = "c12348ce28de40eed0136aa2b644d0ee0650e56c",
)
