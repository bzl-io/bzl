HEAD := $(shell git rev-parse HEAD)

build:
	bazel build //:bzl --workspace_status_command=tools/get_workspace_status.sh

install: build
	cp -f bazel-bin/linux_amd64_stripped/bzl ~/bin/bazel

release: install
	bazel release \
	--owner bzl-io \
	--repo bzl \
	--tag v0.1.5 \
	--notes RELEASE.md \
	--commit $(HEAD) \
	//:bzl
