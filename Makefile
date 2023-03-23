# --------------------------------------------------------------------------
# Makefile for the Fantom Aida World State Manager CLI
#
# v1.0 (2022/09/22) - Initial version
#
# (c) Fantom Foundation, 2022
# --------------------------------------------------------------------------

# what are we building
PROJECT := $(shell basename "$(PWD)")
GO_BIN := $(CURDIR)/build

# compile time variables will be injected into the app
APP_VERSION := 1.0
BUILD_DATE := $(shell date "+%a, %d %b %Y %T")
BUILD_COMPILER := $(shell go version)
BUILD_COMMIT := $(shell git show --format="%H" --no-patch)
BUILD_COMMIT_TIME := $(shell git show --format="%cD" --no-patch)
GOPROXY ?= "https://proxy.golang.org,direct"

.PHONY: all clean help test

all: aida-api-replay aida-worldstate aida-updateset aida-dbmerger aida-trace aida-runarchive aida-runvm aida-stochastic aida-substate

aida-api-replay:
	@cd carmen/go/lib ; \
	./build_libcarmen.sh ; \
	cd ../../.. ; \
	GOPROXY=$(GOPROXY) \
	GOPRIVATE=github.com/Fantom-foundation/Carmen,github.com/Fantom-foundation/go-opera-fvm \
	go build -ldflags "-s -w -X 'github.com/Fantom-foundation/Aida/utils.GitCommit=$(BUILD_COMMIT)'" \
	-o $(GO_BIN)/aida-apireplay \
	./cmd/api-replay-cli

aida-worldstate:
	@go build \
		-ldflags="-X 'github.com/Fantom-foundation/Aida/cmd/worldstate-cli/version.Version=$(APP_VERSION)' -X 'github.com/Fantom-foundation/Aida/cmd/worldstate-cli/version.Time=$(BUILD_DATE)' -X 'github.com/Fantom-foundation/Aida/cmd/worldstate-cli/version.Compiler=$(BUILD_COMPILER)' -X 'github.com/Fantom-foundation/Aida/cmd/worldstate-cli/version.Commit=$(BUILD_COMMIT)' -X 'github.com/Fantom-foundation/Aida/cmd/worldstate-cli/version.CommitTime=$(BUILD_COMMIT_TIME)'" \
		-o $(GO_BIN)/aida-worldstate \
		-v \
		./cmd/worldstate-cli

aida-stochastic:
	@cd carmen/go/lib ; \
	./build_libcarmen.sh ; \
	cd ../../.. ; \
	GOPROXY=$(GOPROXY) \
	GOPRIVATE=github.com/Fantom-foundation/Carmen,github.com/Fantom-foundation/go-opera-fvm \
	go build -ldflags "-s -w -X 'github.com/Fantom-foundation/Aida/utils.GitCommit=$(BUILD_COMMIT)'" \
       	-o $(GO_BIN)/aida-stochastic \
	./cmd/stochastic-cli

aida-trace:
	@cd carmen/go/lib ; \
	./build_libcarmen.sh ; \
	cd ../../.. ; \
	GOPROXY=$(GOPROXY) \
	GOPRIVATE=github.com/Fantom-foundation/Carmen,github.com/Fantom-foundation/go-opera-fvm \
	go build -ldflags "-s -w -X 'github.com/Fantom-foundation/Aida/utils.GitCommit=$(BUILD_COMMIT)'" \
	-o $(GO_BIN)/aida-trace \
	./cmd/trace-cli



aida-runarchive:
	@cd carmen/go/lib ; \
	./build_libcarmen.sh ; \
	cd ../../.. ; \
	GOPROXY=$(GOPROXY) \
	GOPRIVATE=github.com/Fantom-foundation/Carmen,github.com/Fantom-foundation/go-opera-fvm \
	go build -ldflags "-s -w -X 'github.com/Fantom-foundation/Aida/utils.GitCommit=$(BUILD_COMMIT)'" \
	-o $(GO_BIN)/aida-runarchive \
	./cmd/runarchive-cli

aida-runvm:
	@cd carmen/go/lib ; \
	./build_libcarmen.sh ; \
	cd ../../.. ; \
	GOPROXY=$(GOPROXY) \
	GOPRIVATE=github.com/Fantom-foundation/Carmen,github.com/Fantom-foundation/go-opera-fvm \
	go build -ldflags "-s -w -X 'github.com/Fantom-foundation/Aida/utils.GitCommit=$(BUILD_COMMIT)'" \
	-o $(GO_BIN)/aida-runvm \
	./cmd/runvm-cli

aida-substate:
	@cd carmen/go/lib ; \
	./build_libcarmen.sh ; \
	cd ../../.. ; \
	GOPROXY=$(GOPROXY) \
	GOPRIVATE=github.com/Fantom-foundation/Carmen,github.com/Fantom-foundation/go-opera-fvm \
	go build -ldflags "-s -w -X 'github.com/Fantom-foundation/Aida/utils.GitCommit=$(BUILD_COMMIT)'" \
       	-o $(GO_BIN)/aida-substate \
	./cmd/substate-cli

aida-updateset:
	@cd carmen/go/lib ; \
	./build_libcarmen.sh ; \
	cd ../../.. ; \
	GOPROXY=$(GOPROXY) \
	go build -ldflags "-s -w" \
	-o $(GO_BIN)/aida-updateset \
	./cmd/updateset-cli

aida-dbmerger:
	@cd carmen/go/lib ; \
	./build_libcarmen.sh ; \
	cd ../../.. ; \
	GOPROXY=$(GOPROXY) \
	go build -ldflags "-s -w" \
	-o $(GO_BIN)/aida-dbmerger \
	./cmd/db-merger

test:
	@go test ./...

clean:
	cd carmen/go ; \
	rm -f lib/libstate.so ; \
	cd ../cpp ; \
	bazel clean ; \
	cd ../.. ; \
	rm -fr ./build/*

help: Makefile
	@echo "Choose a make command in "$(PROJECT)":"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo
