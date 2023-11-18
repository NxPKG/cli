NXPKG_VERSION = $(shell cat ../version.txt | sed -n '1 p')
NXPKG_TAG = $(shell cat ../version.txt | sed -n '2 p')

EXT :=
ifeq ($(OS),Windows_NT)
	UNAME := Windows
	EXT = .exe
else
	UNAME := $(shell uname -s)
endif

GOARCH:=$(shell go env GOARCH | xargs)
GOOS:=$(shell go env GOOS | xargs)

# Strip debug info
GO_FLAGS += "-ldflags=-s -w"

# Avoid embedding the build path in the executable for more reproducible builds
GO_FLAGS += -trimpath

CLI_DIR = $(shell pwd)

# allow opting in to the rust codepaths
GO_TAG ?= rust

GO_FILES = $(shell find . -name "*.go")
SRC_FILES = $(shell find . -name "*.go" | grep -v "_test.go")
GENERATED_FILES = internal/nxpkgdprotocol/nxpkgd.pb.go internal/nxpkgdprotocol/nxpkgd_grpc.pb.go

# We do not set go-nxpkg as a dependency because the Rust build.rs
# script will call it for us and copy over the binary
nxpkg:
	cargo build -p nxpkg

nxpkg-prod:
	cargo build --release --manifest-path ../crates/nxpkgrepo/Cargo.toml

nxpkg-capnp:
	cd ../crates/nxpkgrepo-lib/src/hash && capnp compile -I std -ogo proto.capnp && mv ./proto.capnp.go ../../../../cli/internal/fs/hash/capnp

go-nxpkg$(EXT): $(GENERATED_FILES) $(SRC_FILES) go.mod nxpkgrepo-ffi-install
	CGO_ENABLED=1 go build -tags $(GO_TAG) -o go-nxpkg$(EXT) ./cmd/nxpkg


.PHONY: nxpkgrepo-ffi-install
nxpkgrepo-ffi-install: nxpkgrepo-ffi nxpkgrepo-ffi-copy-bindings
	cp ../crates/nxpkgrepo-ffi/target/debug/libnxpkgrepo_ffi.a ./internal/ffi/libnxpkgrepo_ffi_$(GOOS)_$(GOARCH).a

.PHONY: nxpkgrepo-ffi
nxpkgrepo-ffi:
	cd ../crates/nxpkgrepo-ffi && cargo build --target-dir ./target

.PHONY: nxpkgrepo-ffi-copy-bindings
nxpkgrepo-ffi-copy-bindings:
	cp ../crates/nxpkgrepo-ffi/bindings.h ./internal/ffi/bindings.h

#
# ffi cross compiling
#
# these targets are used to build the ffi library for each platform
# when doing a release. they _may_ work on your local machine, but
# they're not intended to be used for development.
#

.PHONY: nxpkgrepo-ffi-install-windows-amd64
nxpkgrepo-ffi-install-windows-amd64: nxpkgrepo-ffi-windows-amd64 nxpkgrepo-ffi-copy-bindings
	cp ../crates/nxpkgrepo-ffi/target/x86_64-pc-windows-gnu/release/libnxpkgrepo_ffi.a ./internal/ffi/libnxpkgrepo_ffi_windows_amd64.a

.PHONY: nxpkgrepo-ffi-install-darwin-arm64
nxpkgrepo-ffi-install-darwin-arm64: nxpkgrepo-ffi-darwin-arm64 nxpkgrepo-ffi-copy-bindings
	cp ../crates/nxpkgrepo-ffi/target/aarch64-apple-darwin/release/libnxpkgrepo_ffi.a ./internal/ffi/libnxpkgrepo_ffi_darwin_arm64.a

.PHONY: nxpkgrepo-ffi-install-darwin-amd64
nxpkgrepo-ffi-install-darwin-amd64: nxpkgrepo-ffi-darwin-amd64 nxpkgrepo-ffi-copy-bindings
	cp ../crates/nxpkgrepo-ffi/target/x86_64-apple-darwin/release/libnxpkgrepo_ffi.a ./internal/ffi/libnxpkgrepo_ffi_darwin_amd64.a

.PHONY: nxpkgrepo-ffi-install-linux-arm64
nxpkgrepo-ffi-install-linux-arm64: nxpkgrepo-ffi-linux-arm64 nxpkgrepo-ffi-copy-bindings
	cp ../crates/nxpkgrepo-ffi/target/aarch64-unknown-linux-musl/release/libnxpkgrepo_ffi.a ./internal/ffi/libnxpkgrepo_ffi_linux_arm64.a

.PHONY: nxpkgrepo-ffi-install-linux-amd64
nxpkgrepo-ffi-install-linux-amd64: nxpkgrepo-ffi-linux-amd64 nxpkgrepo-ffi-copy-bindings
	cp ../crates/nxpkgrepo-ffi/target/x86_64-unknown-linux-musl/release/libnxpkgrepo_ffi.a ./internal/ffi/libnxpkgrepo_ffi_linux_amd64.a

.PHONY: nxpkgrepo-ffi-windows-amd64
nxpkgrepo-ffi-windows-amd64:
	cd ../crates/nxpkgrepo-ffi && cargo build --release --target-dir ./target --target x86_64-pc-windows-gnu

.PHONY: nxpkgrepo-ffi-darwin-arm64
nxpkgrepo-ffi-darwin-arm64:
	cd ../crates/nxpkgrepo-ffi && cargo build --release --target-dir ./target --target aarch64-apple-darwin

.PHONY: nxpkgrepo-ffi-darwin-amd64
nxpkgrepo-ffi-darwin-amd64:
	cd ../crates/nxpkgrepo-ffi && cargo build --release --target-dir ./target --target x86_64-apple-darwin

.PHONY: nxpkgrepo-ffi-linux-arm64
nxpkgrepo-ffi-linux-arm64:
	cd ../crates/nxpkgrepo-ffi && CC="zig cc -target aarch64-linux-musl" cargo build --release --target-dir ./target --target aarch64-unknown-linux-musl

.PHONY: nxpkgrepo-ffi-linux-amd64
nxpkgrepo-ffi-linux-amd64:
	cd ../crates/nxpkgrepo-ffi && CC="zig cc -target x86_64-linux-musl" cargo build --release --target-dir ./target --target x86_64-unknown-linux-musl

#
# end
#

.PHONY: nxpkgrepo-ffi-proto
nxpkgrepo-ffi-proto:
	protoc -I../crates/ ../crates/nxpkgrepo-ffi/messages.proto --go_out=./internal/

protoc: internal/nxpkgdprotocol/nxpkgd.proto
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		internal/nxpkgdprotocol/nxpkgd.proto

$(GENERATED_FILES): internal/nxpkgdprotocol/nxpkgd.proto
	make protoc

compile-protos: $(GENERATED_FILES)

ewatch: scripts/...
	nodemon --exec "make e2e" -e .ts,.go

check-go-version:
	@go version | grep ' go1\.18\.0 ' || (echo 'Please install Go version 1.18.0' && false)

# This "NXPKG_RACE" variable exists at the request of a user on GitHub who
# wants to run "make test-go" on an unsupported version of macOS (version 10.9).
# Go's race detector does not run correctly on that version. With this flag
# you can run "NXPKG_RACE= make test-go" to disable the race detector.
NXPKG_RACE ?= -race

ifeq ($(UNAME), Windows)
	NXPKG_RACE=
endif

clean-go:
	go clean -testcache -r

test-go: $(GENERATED_FILES) $(GO_FILES) go.mod go.sum nxpkgrepo-ffi-install
	go test $(NXPKG_RACE) -tags $(GO_TAG) ./...

# protos need to be compiled before linting, since linting needs to pick up
# some types from the generated code
lint-go: $(GENERATED_FILES) $(GO_FILES) go.mod go.sum
	golangci-lint run --new-from-rev=main

fmt-go: $(GO_FILES) go.mod go.sum
	go fmt ./...

install: | ./package.json
	pnpm install --filter=cli

corepack:
	which corepack || npm install -g corepack@latest
	corepack enable

e2e: corepack install nxpkg
	node -r esbuild-register scripts/e2e/e2e.ts

# Expects nxpkg to be built and up to date
# Only should be used by CI
e2e-prebuilt: corepack install
	node -r esbuild-register scripts/e2e/e2e.ts

cmd/nxpkg/version.go: ../version.txt
	# Update this atomically to avoid issues with this being overwritten during use
	node -e 'console.log(`package main\n\nconst nxpkgVersion = "$(NXPKG_VERSION)"`)' > cmd/nxpkg/version.go.txt
	mv cmd/nxpkg/version.go.txt cmd/nxpkg/version.go

build: install
	cd $(CLI_DIR)/../ && pnpm build:nxpkg
	cd $(CLI_DIR)/../ && pnpm install --filter=create-nxpkg && pnpm nxpkg-prebuilt build --filter=create-nxpkg...
	cd $(CLI_DIR)/../ && pnpm install --filter=@nxpkg/codemod && pnpm nxpkg-prebuilt build --filter=@nxpkg/codemod...
	cd $(CLI_DIR)/../ && pnpm install --filter=nxpkg-ignore && pnpm nxpkg-prebuilt build --filter=nxpkg-ignore...
	cd $(CLI_DIR)/../ && pnpm install --filter=@nxpkg/workspaces && pnpm nxpkg-prebuilt build --filter=@nxpkg/workspaces...
	cd $(CLI_DIR)/../ && pnpm install --filter=@nxpkg/gen && pnpm nxpkg-prebuilt build --filter=@nxpkg/gen...
	cd $(CLI_DIR)/../ && pnpm install --filter=eslint-plugin-nxpkg && pnpm nxpkg-prebuilt build --filter=eslint-plugin-nxpkg...
	cd $(CLI_DIR)/../ && pnpm install --filter=eslint-config-nxpkg && pnpm nxpkg-prebuilt build --filter=eslint-config-nxpkg...
	cd $(CLI_DIR)/../ && pnpm install --filter=@nxpkg/types && pnpm nxpkg-prebuilt build --filter=@nxpkg/types...

.PHONY: prepublish
prepublish: compile-protos cmd/nxpkg/version.go
	make -j3 bench/nxpkg test-go

.PHONY: publish-nxpkg-cross
publish-nxpkg-cross: prepublish
	goreleaser release --rm-dist -f cross-release.yml

.PHONY: publish-nxpkg-darwin
publish-nxpkg-darwin: prepublish
	goreleaser release --rm-dist -f darwin-release.yml

.PHONY: snapshot-nxpkg-cross
snapshot-nxpkg-cross:
	goreleaser release --snapshot --rm-dist -f cross-release.yml

.PHONY: snapshot-nxpkg-darwin
snapshot-nxpkg-darwin:
	goreleaser release --snapshot --rm-dist -f darwin-release.yml

.PHONY: snapshot-lib-nxpkg-darwin
snapshot-lib-nxpkg-darwin:
	goreleaser release --snapshot --rm-dist -f darwin-lib.yml

.PHONY: snapshot-lib-nxpkg-cross
snapshot-lib-nxpkg-cross:
	goreleaser release --snapshot --rm-dist -f cross-lib.yml

.PHONY: build-lib-nxpkg-darwin
build-lib-nxpkg-darwin:
	goreleaser release --rm-dist -f darwin-lib.yml

.PHONY: build-go-nxpkg-darwin
build-go-nxpkg-darwin:
	goreleaser release --rm-dist -f darwin-release.yml

.PHONY: build-go-nxpkg-cross
build-go-nxpkg-cross:
	goreleaser release --rm-dist -f cross-release.yml

.PHONY: build-lib-nxpkg-cross
build-lib-nxpkg-cross:
	goreleaser release --rm-dist -f cross-lib.yml

.PHONY: stage-release
stage-release: cmd/nxpkg/version.go
	echo "Version: $(NXPKG_VERSION)"
	echo "Tag: $(NXPKG_TAG)"
	cat $(CLI_DIR)/../version.txt
	git diff -- $(CLI_DIR)/../version.txt
	git status
	@test "" = "`git cherry`" || (echo "Refusing to publish with unpushed commits" && false)

	# Stop if versions are not updated.
	@test "" != "`git diff -- $(CLI_DIR)/../version.txt`" || (echo "Refusing to publish with unupdated version.txt" && false)
	@test "" != "`git diff -- $(CLI_DIR)/cmd/nxpkg/version.go`" || (echo "Refusing to publish with unupdated version.go" && false)

	# Prepare the packages.
	cd $(CLI_DIR)/../packages/nxpkg && pnpm version "$(NXPKG_VERSION)" --allow-same-version
	cd $(CLI_DIR)/../packages/create-nxpkg && pnpm version "$(NXPKG_VERSION)" --allow-same-version
	cd $(CLI_DIR)/../packages/nxpkg-codemod && pnpm version "$(NXPKG_VERSION)" --allow-same-version
	cd $(CLI_DIR)/../packages/nxpkg-ignore && pnpm version "$(NXPKG_VERSION)" --allow-same-version
	cd $(CLI_DIR)/../packages/nxpkg-workspaces && pnpm version "$(NXPKG_VERSION)" --allow-same-version
	cd $(CLI_DIR)/../packages/nxpkg-gen && pnpm version "$(NXPKG_VERSION)" --allow-same-version
	cd $(CLI_DIR)/../packages/eslint-plugin-nxpkg && pnpm version "$(NXPKG_VERSION)" --allow-same-version
	cd $(CLI_DIR)/../packages/eslint-config-nxpkg && pnpm version "$(NXPKG_VERSION)" --allow-same-version
	cd $(CLI_DIR)/../packages/nxpkg-types && pnpm version "$(NXPKG_VERSION)" --allow-same-version

	git checkout -b staging-$(NXPKG_VERSION)
	git commit -anm "publish $(NXPKG_VERSION) to registry"
	git tag "v$(NXPKG_VERSION)"
	git push origin staging-$(NXPKG_VERSION) --tags --force

.PHONY: publish-nxpkg
publish-nxpkg: clean build
	echo "Version: $(NXPKG_VERSION)"
	echo "Tag: $(NXPKG_TAG)"

	# Include the patch in the log.
	git format-patch HEAD~1 --stdout | cat

	npm config set --location=project "//registry.npmjs.org/:_authToken" $(NPM_TOKEN)

	# Publishes the native npm modules.
	goreleaser release --rm-dist -f combined-shim.yml $(SKIP_PUBLISH)

	# Split packing from the publish step so that npm locates the correct .npmrc file.
	cd $(CLI_DIR)/../packages/nxpkg && pnpm pack --pack-destination=$(CLI_DIR)/../
	cd $(CLI_DIR)/../packages/create-nxpkg && pnpm pack --pack-destination=$(CLI_DIR)/../
	cd $(CLI_DIR)/../packages/nxpkg-codemod && pnpm pack --pack-destination=$(CLI_DIR)/../
	cd $(CLI_DIR)/../packages/nxpkg-ignore && pnpm pack --pack-destination=$(CLI_DIR)/../
	cd $(CLI_DIR)/../packages/nxpkg-workspaces && pnpm pack --pack-destination=$(CLI_DIR)/../
	cd $(CLI_DIR)/../packages/nxpkg-gen && pnpm pack --pack-destination=$(CLI_DIR)/../
	cd $(CLI_DIR)/../packages/eslint-plugin-nxpkg && pnpm pack --pack-destination=$(CLI_DIR)/../
	cd $(CLI_DIR)/../packages/eslint-config-nxpkg && pnpm pack --pack-destination=$(CLI_DIR)/../
	cd $(CLI_DIR)/../packages/nxpkg-types && pnpm pack --pack-destination=$(CLI_DIR)/../

ifneq ($(SKIP_PUBLISH),--skip-publish)
	# Publish the remaining JS packages in order to avoid race conditions.
	cd $(CLI_DIR)/../
	npm publish -ddd --tag $(NXPKG_TAG) $(CLI_DIR)/../nxpkg-$(NXPKG_VERSION).tgz
	npm publish -ddd --tag $(NXPKG_TAG) $(CLI_DIR)/../create-nxpkg-$(NXPKG_VERSION).tgz
	npm publish -ddd --tag $(NXPKG_TAG) $(CLI_DIR)/../nxpkg-codemod-$(NXPKG_VERSION).tgz
	npm publish -ddd --tag $(NXPKG_TAG) $(CLI_DIR)/../nxpkg-ignore-$(NXPKG_VERSION).tgz
	npm publish -ddd --tag $(NXPKG_TAG) $(CLI_DIR)/../nxpkg-workspaces-$(NXPKG_VERSION).tgz
	npm publish -ddd --tag $(NXPKG_TAG) $(CLI_DIR)/../nxpkg-gen-$(NXPKG_VERSION).tgz
	npm publish -ddd --tag $(NXPKG_TAG) $(CLI_DIR)/../eslint-plugin-nxpkg-$(NXPKG_VERSION).tgz
	npm publish -ddd --tag $(NXPKG_TAG) $(CLI_DIR)/../eslint-config-nxpkg-$(NXPKG_VERSION).tgz
	npm publish -ddd --tag $(NXPKG_TAG) $(CLI_DIR)/../nxpkg-types-$(NXPKG_VERSION).tgz
endif

demo/lage: install
	node $(CLI_DIR)/scripts/generate.mjs lage

demo/lerna: install
	node $(CLI_DIR)/scripts/generate.mjs lerna

demo/nx: install
	node $(CLI_DIR)/scripts/generate.mjs nx

demo/nxpkg: install
	node $(CLI_DIR)/scripts/generate.mjs nxpkg

demo: demo/lage demo/lerna demo/nx demo/nxpkg

bench/lerna: demo/lerna
	cd $(CLI_DIR)/demo/lerna && node_modules/.bin/lerna run build

bench/lage: demo/lage
	cd $(CLI_DIR)/demo/lage && node_modules/.bin/lage build

bench/nx: demo/nx
	cd $(CLI_DIR)/demo/nx && node_modules/.bin/nx run-many --target=build --all

bench/nxpkg: demo/nxpkg nxpkg
	cd $(CLI_DIR)/demo/nxpkg && $(CLI_DIR)/nxpkg run test

bench: bench/lerna bench/lage bench/nx bench/nxpkg

clean: clean-go clean-build clean-demo clean-rust

clean-rust:
	cargo clean

clean-build:
	rm -f nxpkg

clean-demo:
	rm -rf node_modules
	rm -rf demo

# use target fixture-<some directory under integration_tests/_fixtures> to set up the testbed directory
.PHONY=fixture-%
fixture-%:
	$(eval $@_FIXTURE := $(@:fixture-%=%))
	@echo "fixture setup $($@_FIXTURE)"
	rm -rf testbed
	mkdir -p testbed
	../nxpkgrepo-tests/integration/tests/_helpers/setup_monorepo.sh ./testbed $($@_FIXTURE)