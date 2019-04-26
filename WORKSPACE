workspace(name = "io_bzl_bzl")

load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

http_archive(
    name = "build_stack_rules_proto",
    sha256 = "06bd105c1b0f8ea7c2827da045fcb83f44dd2f78e03d98abc1de4ec21e45c9d6",
    strip_prefix = "rules_proto-218e598f693964014fc9d3dbc2dfb986fbe09d81",
    urls = ["https://github.com/stackb/rules_proto/archive/218e598f693964014fc9d3dbc2dfb986fbe09d81.tar.gz"],
)

load(
    "@build_stack_rules_proto//:deps.bzl",
    "bazel_gazelle",
    "com_github_bazelbuild_buildtools",
    "io_bazel_rules_go",
)

io_bazel_rules_go()

bazel_gazelle()

com_github_bazelbuild_buildtools()

load("@io_bazel_rules_go//go:def.bzl", "go_register_toolchains", "go_rules_dependencies")

go_rules_dependencies()

go_register_toolchains()

load("@bazel_gazelle//:deps.bzl", "gazelle_dependencies", "go_repository")

gazelle_dependencies()

# gazelle:repo bazel_gazelle

# ================================================================

go_repository(
    name = "com_github_google_go_github",
    commit = "61bf8a6a2009b7fb88173ed454c77fb8013b8852",
    importpath = "github.com/google/go-github",
)

go_repository(
    name = "com_github_google_go_querystring",
    commit = "53e6ce116135b80d037921a7fdd5138cf32d7a8a",
    importpath = "github.com/google/go-querystring",
)

go_repository(
    name = "com_github_davecgh_go_spew",
    commit = "a476722483882dd40b8111f0eb64e1d7f43f56e4",
    importpath = "github.com/davecgh/go-spew",
)

go_repository(
    name = "com_golang_google_genproto",
    commit = "3273178ea4684acc4f512f7bef7349dd72db88f6",
    importpath = "google.golang.org/genproto",
)

go_repository(
    name = "com_github_gregjones_httpcache",
    commit = "22a0b1feae53974ed4cfe27bcce70dba061cc5fd",
    importpath = "github.com/gregjones/httpcache",
)

go_repository(
    name = "com_github_peterbourgon_diskv",
    commit = "53ef9e43a0bc608e737e6bfed35207ad9cb1ad54",
    importpath = "github.com/peterbourgon/diskv",
)

go_repository(
    name = "com_github_google_btree",
    commit = "316fb6d3f031ae8f4d457c6c5186b9e3ded70435",
    importpath = "github.com/google/btree",
)

go_repository(
    name = "com_github_urfave_cli",
    commit = "44cb242eeb4d76cc813fdc69ba5c4b224677e799",
    importpath = "github.com/urfave/cli",
)

go_repository(
    name = "com_github_dustin_go_humanize",
    commit = "6d15c0ae71e55ed645c21ac4945aaadbc0e9a590",
    importpath = "github.com/dustin/go-humanize",
)

go_repository(
    name = "com_github_cheggaaa_pb",
    commit = "657164d0228d6bebe316fdf725c69f131a50fb10",
    importpath = "github.com/cheggaaa/pb",
)

go_repository(
    name = "com_github_pkg_errors",
    commit = "f15c970de5b76fac0b59abb32d62c17cc7bed265",
    importpath = "github.com/pkg/errors",
)

go_repository(
    name = "org_golang_x_tools",
    commit = "9b61fcc4c548d69663d915801fc4b42a43b6cd9c",
    importpath = "github.com/golang/tools",
)

go_repository(
    name = "org_golang_x_sync",
    commit = "fd80eb99c8f653c847d294a001bdf2a3a6f768f5",
    importpath = "github.com/golang/sync",
)

go_repository(
    name = "com_google_bazelbuild_buildtools",
    commit = "8135d8f1de24e6cb453be56ce061c922bf279f2c",
    importpath = "github.com/bazelbuild/buildtools",
)

go_repository(
    name = "com_github_mitchellh_go_homedir",
    commit = "b8bc1bf767474819792c23f32d8286a45736f1c6",
    importpath = "github.com/mitchellh/go-homedir",
)

# There exist multiple options for terminal-based progress bars, but this is the
# only one I could find that cross-compiles cleanly with current version of
# rules_go.  Not super fancy though.
go_repository(
    name = "com_github_mitchellh_ioprogress",
    commit = "8163955264568045f462ae7e2d6d07b2001fc997",
    importpath = "github.com/mitchellh/ioprogress",
)

go_repository(
    name = "com_github_joeybloggs_go_download",
    commit = "26df310821f0e7614a736a2ee38ecd7e2f6ef6da",
    importpath = "github.com/joeybloggs/go-download",
)

go_repository(
    name = "com_github_rs_cors",
    commit = "eabcc6af4bbe5ad3a949d36450326a2b0b9894b8",  # Aug 1
    importpath = "github.com/rs/cors",
)

go_repository(
    name = "com_github_gorilla_mux",
    commit = "3f19343c7d9ce75569b952758bd236af94956061",
    importpath = "github.com/gorilla/mux",
)

go_repository(
    name = "com_github_gorilla_context",
    commit = "08b5f424b9271eedf6f9f0ce86cb9396ed337a42",
    importpath = "github.com/gorilla/context",
)

go_repository(
    name = "com_github_vbauerster_mpb",
    commit = "d3da256ab98a80013319df621e48db638748c044",
    importpath = "github.com/vbauerster/mpb",
)

go_repository(
    name = "com_github_matttproud_golang_protobuf_extensions",
    commit = "c12348ce28de40eed0136aa2b644d0ee0650e56c",
    importpath = "github.com/matttproud/golang_protobuf_extensions",
)

go_repository(
    name = "com_github_fatih_color",
    commit = "3f9d52f7176a6927daacff70a3e8d1dc2025c53e",
    importpath = "github.com/fatih/color",
)
