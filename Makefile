
PROJECT = goutil

ifeq "$(GOPATH)" ""
  $(error Please set the environment variable GOPATH before running `make`)
endif

FAIL_ON_STDOUT := awk '{ print } END { if (NR > 0) { exit 1 } }'

CURDIR := $(shell pwd)
path_to_add := $(addsuffix /bin,$(subst :,/bin:,$(GOPATH)))
export PATH := $(path_to_add):$(PATH)

GO        := GO111MODULE=off go
GOBUILD   := CGO_ENABLED=0 $(GO) build $(BUILD_FLAG)
GOTEST    := CGO_ENABLED=1 $(GO) test -p 3
OVERALLS  := CGO_ENABLED=1 overalls
GOVERALLS := goveralls

ARCH      := "`uname -s`"
LINUX     := "Linux"
MAC       := "Darwin"
PACKAGE_LIST  := go list ./...
PACKAGES  := $$($(PACKAGE_LIST))
PACKAGE_DIRECTORIES := $(PACKAGE_LIST) | sed 's|$(PROJECT)/||'
FILES     := $$(find $$($(PACKAGE_DIRECTORIES)) -name "*.go")

GOFAIL_ENABLE  := $$(find $$PWD/ -type d | grep -vE "(\.git|_tools)" | xargs gofail enable)
GOFAIL_DISABLE := $$(find $$PWD/ -type d | grep -vE "(\.git|_tools)" | xargs gofail disable)

LDFLAGS += -X "github.com/zhengcf/goutil/printer.AppReleaseVersion=$(shell git describe --tags --dirty)"
LDFLAGS += -X "github.com/zhengcf/goutil/printer.AppBuildTS=$(shell date -u '+%Y-%m-%d %I:%M:%S')"
LDFLAGS += -X "github.com/zhengcf/goutil/printer.AppGitHash=$(shell git rev-parse HEAD)"
LDFLAGS += -X "github.com/zhengcf/goutil/printer.AppGitBranch=$(shell git rev-parse --abbrev-ref HEAD)"
LDFLAGS += -X "github.com/zhengcf/goutil/printer.GoVersion=$(shell go version)"

CHECK_LDFLAGS += $(LDFLAGS)

TARGET = ""

.PHONY: all build update clean check

default: check all buildsucc

buildsucc:
	@echo Build ${PROJECT} successfully!

all: build

build:
	$(GOBUILD)  -ldflags '$(LDFLAGS)' -o server/${PROJECT} server/${PROJECT}.go

# The retool tools.json is setup from hack/retool-install.sh
check-setup:
	@which retool >/dev/null 2>&1 || go get github.com/twitchtv/retool
	@GO111MODULE=off retool sync

check: check-setup fmt #lint

# These need to be fixed before they can be ran regularly
check-fail: goword check-static check-slow

fmt:
	@echo "gofmt (simplify)"
	@gofmt -s -l -w $(FILES) 2>&1 | $(FAIL_ON_STDOUT)

lint:
	@echo "linting"
	@CGO_ENABLED=0 retool do revive -formatter friendly -config revive.toml $(PACKAGES)

goword:
	retool do goword $(FILES) 2>&1 | $(FAIL_ON_STDOUT)

check-static:
	@ # vet and fmt have problems with vendor when ran through metalinter
	CGO_ENABLED=0 retool do gometalinter.v2 --disable-all --deadline 120s \
	  --enable misspell \
	  --enable megacheck \
	  --enable ineffassign \
	  $$($(PACKAGE_DIRECTORIES))

check-slow:
	CGO_ENABLED=0 retool do gometalinter.v2 --disable-all \
	  --enable errcheck \
	  $$($(PACKAGE_DIRECTORIES))
	CGO_ENABLED=0 retool do gosec $$($(PACKAGE_DIRECTORIES))


clean:
	$(GO) clean -i ./...
