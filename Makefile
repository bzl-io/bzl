HEAD := $(shell git rev-parse HEAD)

install:
	bazel build //:bzl --workspace_status_command=tools/get_workspace_status.sh && cp --f bazel-bin/linux_amd64_stripped/bzl ~/bin/bazel

release:
	bzl release \
	--platform linux_amd64 \
	--platform linux_386 \
	--platform windows_amd64 \
	--platform windows_386 \
	--owner bzl-io \
	--repo bzl \
	--tag v0.1.2 \
	--notes RELEASE_NOTES.md \
	--commit $(HEAD) \
	//go/bzl:bzl
